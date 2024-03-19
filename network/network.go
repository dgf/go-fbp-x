package network

import (
	"fmt"

	"github.com/dgf/go-fbp-x/process"
)

type packet struct {
	target chan<- any
	data   string
	Connection
}

type stream struct {
	source <-chan any
	target chan<- any
	Connection
}

type Network struct {
	factory    map[string]func() process.Process
	processes  map[string]process.Process
	initialIPs []packet
	streams    []stream
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

func (n *Network) connect(connections []Connection) error {
	n.streams = []stream{}

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
		} else {
			n.streams = append(n.streams, stream{output.Channel, input.Channel, c})
		}
	}

	return nil
}

func Create(graph Graph, out chan<- string) (*Network, error) {
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

func (n *Network) Run(traces chan<- Trace) {
	for _, s := range n.streams {
		go func(connection Connection, source <-chan any, target chan<- any) {
			for value := range source {
				traces <- Trace{value, connection}
				target <- value
			}
		}(s.Connection, s.source, s.target)
	}

	for _, p := range n.initialIPs {
		go func(connection Connection, data string, target chan<- any) {
			traces <- Trace{data, connection}
			target <- data
		}(p.Connection, p.data, p.target)
	}
}
