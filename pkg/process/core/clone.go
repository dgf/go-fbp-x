package core

import "github.com/dgf/go-fbp-x/pkg/process"

type clone struct {
	in  chan any
	out chan any
}

func Clone() process.Process {
	c := &clone{
		in:  make(chan any, 1),
		out: make(chan any, 1),
	}

	go func() {
		defer close(c.out)

		for in := range c.in {
			c.out <- in
		}
	}()

	return c
}

func (*clone) Description() string {
	return "Clone input, e.g. useful to mux and demux packets."
}

func (c *clone) Inputs() map[string]process.Input {
	return map[string]process.Input{
		"in": {Channel: c.in, IPType: process.AnyIP},
	}
}

func (c *clone) Outputs() map[string]process.Output {
	return map[string]process.Output{
		"out": {Channel: c.out, IPType: process.AnyIP},
	}
}
