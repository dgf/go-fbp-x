package network

import (
	"context"
	"fmt"
	"sync"

	"github.com/dgf/go-fbp-x/pkg/dsl"
	"github.com/dgf/go-fbp-x/pkg/process"
)

type Network interface {
	Run(ctx context.Context, graph dsl.Graph, traces chan<- Trace) error
}

type packet struct {
	target chan<- any
	data   string
	dsl.Connection
}

type targetChannel struct {
	channel chan<- any
	dsl.Link
}

type demux struct {
	source  <-chan any
	targets []targetChannel
}

type net struct {
	factory    process.Factory
	processes  map[string]process.Process
	demuxes    map[dsl.Link]*demux
	initialIPs []packet
}

func NewNetwork(factory process.Factory) Network {
	return &net{factory: factory}
}

func (n *net) reference(components map[string]dsl.Process) error {
	n.processes = map[string]process.Process{}

	for component, process := range components {
		if pf, ok := n.factory.Create(process.Name); !ok {
			return fmt.Errorf("process %q not available", process)
		} else {
			n.processes[component] = pf
		}
	}

	return nil
}

func (n *net) connect(connections []dsl.Connection) error {
	n.demuxes = map[dsl.Link]*demux{}

	for _, c := range connections {
		if target, ok := n.processes[c.Target.Component]; !ok {
			return fmt.Errorf("target %q not registered", c.Target.Component)
		} else if input, ok := target.Inputs()[c.Target.Port]; !ok {
			return fmt.Errorf("input %q on target %q not available", c.Target.Port, c.Target.Component)
		} else if len(c.Data) > 0 {
			n.initialIPs = append(n.initialIPs, packet{input.Channel, c.Data, c})
		} else if source, ok := n.processes[c.Source.Component]; !ok {
			return fmt.Errorf("source %q not registered", c.Source.Component)
		} else if output, ok := source.Outputs()[c.Source.Port]; !ok {
			return fmt.Errorf("output %q on source %q not available", c.Source.Port, c.Source.Component)
		} else if !process.IsCompatibleIPType(output.IPType, input.IPType) {
			return fmt.Errorf("unmatched connection type from %v to %v", c.Source, c.Target)
		} else if d, ok := n.demuxes[c.Source]; ok {
			d.targets = append(d.targets, targetChannel{input.Channel, c.Target})
		} else {
			n.demuxes[c.Source] = &demux{output.Channel, []targetChannel{{input.Channel, c.Target}}}
		}
	}

	return nil
}

func (n *net) Run(ctx context.Context, graph dsl.Graph, traces chan<- Trace) error {
	wg := sync.WaitGroup{}

	if err := n.reference(graph.Components); err != nil {
		return err
	} else if err := n.connect(graph.Connections); err != nil {
		return err
	}

	wg.Add(len(n.demuxes))
	for l, d := range n.demuxes {
		go func(link dsl.Link, source <-chan any, targets []targetChannel) {
			defer wg.Done()

			for value := range source {
				for _, target := range targets {
					select {
					case <-ctx.Done():
						return
					default:
						traces <- Trace{value, dsl.Connection{Source: link, Target: target.Link}}
						target.channel <- value
					}
				}
			}
		}(l, d.source, d.targets)
	}

	wg.Add(len(n.initialIPs))
	for _, p := range n.initialIPs {
		go func(connection dsl.Connection, data string, target chan<- any) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				return
			default:
				traces <- Trace{data, connection}
				target <- data
			}
		}(p.Connection, p.data, p.target)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		<-ctx.Done()
		for _, p := range n.processes {
			for _, i := range p.Inputs() {
				close(i.Channel)
			}
		}
	}()

	wg.Wait()
	traces <- Trace{"stopped", dsl.Connection{}}
	return nil
}
