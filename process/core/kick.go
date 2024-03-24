package core

import "github.com/dgf/go-fbp-x/process"

type kick struct {
	data chan any
	in   chan any
	out  chan any
}

func Kick() process.Process {
	c := &kick{
		data: make(chan any, 1),
		in:   make(chan any, 1),
		out:  make(chan any, 1),
	}

	go func() {
		d := <-c.data
		for {
			select {
			case d = <-c.data:
				continue
			case <-c.in: // kicks it
				c.out <- d
			}
		}
	}()

	return c
}

func (*kick) Description() string {
	return "Kicks data packet for each input."
}

func (c *kick) Inputs() map[string]process.Input {
	return map[string]process.Input{
		"data": {Channel: c.data, IPType: process.AnyIP},
		"in":   {Channel: c.in, IPType: process.AnyIP},
	}
}

func (c *kick) Outputs() map[string]process.Output {
	return map[string]process.Output{
		"out": {Channel: c.out, IPType: process.AnyIP},
	}
}
