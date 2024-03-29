package core

import "github.com/dgf/go-fbp-x/process"

type kick struct {
	data chan any
	in   chan any
	out  chan any
}

func Kick() process.Process {
	k := &kick{
		data: make(chan any, 1),
		in:   make(chan any, 1),
		out:  make(chan any, 1),
	}

	go func() {
		defer close(k.out)

		d, ok := <-k.data
		if !ok {
			return
		}

		for {
			select {
			case d, ok = <-k.data:
				if !ok {
					return
				}
			case _, ok := <-k.in: // kicks it
				if !ok {
					return
				}
				k.out <- d
			}
		}
	}()

	return k
}

func (*kick) Description() string {
	return "Kicks data packet for each input."
}

func (k *kick) Inputs() map[string]process.Input {
	return map[string]process.Input{
		"data": {Channel: k.data, IPType: process.AnyIP},
		"in":   {Channel: k.in, IPType: process.AnyIP},
	}
}

func (k *kick) Outputs() map[string]process.Output {
	return map[string]process.Output{
		"out": {Channel: k.out, IPType: process.AnyIP},
	}
}
