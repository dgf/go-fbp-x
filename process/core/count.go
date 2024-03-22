package core

import "github.com/dgf/go-fbp-x/process"

type count struct {
	in  chan any
	out chan any
}

func Count() process.Process {
	c := &count{
		in:  make(chan any, 1),
		out: make(chan any, 1),
	}

	go func() {
		cnt := 0
		for range c.in {
			cnt++
			c.out <- cnt
		}
	}()

	return c
}

func (*count) Description() string {
	return "Counts all inputs, outputs an increment for each input."
}

func (c *count) Inputs() map[string]process.Input {
	return map[string]process.Input{
		"in": {Channel: c.in, IPType: process.AnyIP},
	}
}

func (c *count) Outputs() map[string]process.Output {
	return map[string]process.Output{
		"out": {Channel: c.out, IPType: process.NumberIP},
	}
}
