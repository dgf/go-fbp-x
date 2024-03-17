package network

import (
	"fmt"

	"github.com/dgf/go-fbp-x/dsl"
	"github.com/dgf/go-fbp-x/process"
)

type Network struct {
	components map[string]process.Process
	initialIPs []func()
}

func initLibrary(out chan<- string) map[string]func() process.Process {
	processes := map[string]func() process.Process{}

	processes["Counter"] = func() process.Process { return process.Counter() }
	processes["OutputText"] = func() process.Process { return process.OutputText(out) }
	processes["ReadFile"] = func() process.Process { return process.ReadFile() }
	processes["SplitLines"] = func() process.Process { return process.SplitLines() }

	return processes
}

func referenceComponents(graph dsl.Graph, processes map[string]func() process.Process) (map[string]process.Process, error) {
	components := map[string]process.Process{}
	for c, p := range graph.Components {
		if process, ok := processes[p]; !ok {
			return components, fmt.Errorf("process %q not available", p)
		} else {
			components[c] = process()
		}
	}
	return components, nil
}

func Create(graph dsl.Graph, out chan<- string) (*Network, error) {
	network := &Network{}
	processes := initLibrary(out)

	components, err := referenceComponents(graph, processes)
	if err != nil {
		return network, err
	}
	network.components = components

	for _, c := range graph.Connections {
		if target, ok := network.components[c.Target.Component]; !ok {
			return network, fmt.Errorf("target %q not registered", c.Target.Component)
		} else if input, ok := target.Inputs()[c.Target.Port]; !ok {
			return network, fmt.Errorf("input %q on target %q not available", c.Target.Port, c.Target.Component)
		} else {
			if len(c.Data) > 0 { // TODO validate string input
				network.initialIPs = append(network.initialIPs, func() { input.Channel <- c.Data })
			} else if source, ok := network.components[c.Source.Component]; !ok {
				return network, fmt.Errorf("source %q not registered", c.Source.Component)
			} else if output, ok := source.Outputs()[c.Source.Port]; !ok {
				return network, fmt.Errorf("output %q on source %q not available", c.Source.Port, c.Source.Component)
			} else if !process.IsCompatibleIPType(output.IPType, input.IPType) {
				return network, fmt.Errorf("unmatched connection type from %v to %v", c.Source, c.Target)
			} else {
				go func() {
					for value := range output.Channel {
						input.Channel <- value
					}
				}()
			}
		}
	}

	return network, nil
}

func (n *Network) Run() error {
	for _, i := range n.initialIPs {
		go func(init func()) { init() }(i)
	}

	return nil
}
