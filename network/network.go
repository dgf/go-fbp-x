package network

import (
	"fmt"

	"github.com/dgf/go-fbp-x/dsl"
	"github.com/dgf/go-fbp-x/process"
)

type Stream struct {
	source <-chan any
	target chan<- any
}

type Network struct {
	factory    map[string]func() process.Process
	processes  map[string]process.Process
	initialIPs []func()
	streams    []Stream
}

func (n *Network) init(out chan<- string) {
	n.factory = map[string]func() process.Process{}

	n.factory["Counter"] = func() process.Process { return process.Counter() }
	n.factory["OutputText"] = func() process.Process { return process.OutputText(out) }
	n.factory["ReadFile"] = func() process.Process { return process.ReadFile() }
	n.factory["SplitLines"] = func() process.Process { return process.SplitLines() }
}

func (n *Network) reference(components map[string]string) error {
	n.processes = map[string]process.Process{}

	for component, process := range components {
		if factory, ok := n.factory[process]; !ok {
			return fmt.Errorf("process %q not available", process)
		} else {
			n.processes[component] = factory()
		}
	}

	return nil
}

func (n *Network) connect(connections []dsl.Connection) error {
	n.streams = []Stream{}

	for _, c := range connections {
		if target, ok := n.processes[c.Target.Component]; !ok {
			return fmt.Errorf("target %q not registered", c.Target.Component)
		} else if input, ok := target.Inputs()[c.Target.Port]; !ok {
			return fmt.Errorf("input %q on target %q not available", c.Target.Port, c.Target.Component)
		} else if len(c.Data) > 0 {
			n.initialIPs = append(n.initialIPs, func() { input.Channel <- c.Data })
		} else if source, ok := n.processes[c.Source.Component]; !ok {
			return fmt.Errorf("source %q not registered", c.Source.Component)
		} else if output, ok := source.Outputs()[c.Source.Port]; !ok {
			return fmt.Errorf("output %q on source %q not available", c.Source.Port, c.Source.Component)
		} else if !process.IsCompatibleIPType(output.IPType, input.IPType) {
			return fmt.Errorf("unmatched connection type from %v to %v", c.Source, c.Target)
		} else {
			n.streams = append(n.streams, Stream{output.Channel, input.Channel})
		}
	}

	return nil
}

func Create(graph dsl.Graph, out chan<- string) (*Network, error) {
	network := &Network{}
	network.init(out)

	if err := network.reference(graph.Components); err != nil {
		return network, err
	} else if err := network.connect(graph.Connections); err != nil {
		return network, err
	} else {
		return network, nil
	}
}

func (n *Network) Run() error {
	for _, s := range n.streams {
		go func(source <-chan any, target chan<- any) {
			for value := range source {
				target <- value
			}
		}(s.source, s.target)
	}

	for _, i := range n.initialIPs {
		go func(init func()) { init() }(i)
	}

	return nil
}
