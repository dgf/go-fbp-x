package process

import (
	"fmt"

	"github.com/dgf/go-fbp-x/dsl"
)

type Network struct {
	components map[string]Process
	dataInputs []func()
}

type Process interface {
	Inputs() map[string]chan<- string
	Outputs() map[string]<-chan string
}

func initLibrary(out chan<- string) map[string]func() Process {
	processes := map[string]func() Process{}
	processes["Output"] = func() Process { return Output(out) }
	processes["ReadFile"] = func() Process { return ReadFile() }
	return processes
}

func referenceComponents(graph dsl.Graph, processes map[string]func() Process) (map[string]Process, error) {
	components := map[string]Process{}
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
			if len(c.Data) > 0 {
				network.dataInputs = append(network.dataInputs, func() { input <- c.Data })
			} else if source, ok := network.components[c.Source.Component]; !ok {
				return network, fmt.Errorf("source %q not registered", c.Source.Component)
			} else if output, ok := source.Outputs()[c.Source.Port]; !ok {
				return network, fmt.Errorf("output %q on source %q not available", c.Source.Port, c.Source.Component)
			} else {
				go func() {
					for value := range output {
						input <- value
					}
				}()
			}
		}
	}

	return network, nil
}

func (n *Network) Run() error {
	for _, i := range n.dataInputs {
		go func(init func()) { init() }(i)
	}

	return nil
}
