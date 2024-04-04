package core

import (
	"fmt"
	"strconv"

	"github.com/dgf/go-fbp-x/process"
)

type output struct {
	in chan any
}

func Output(out chan<- string) process.Process {
	o := &output{in: make(chan any, 1)}

	go func() {
		for i := range o.in {
			if s, ok := i.(string); ok {
				out <- s
			} else if n, ok := i.(int); ok {
				out <- strconv.Itoa(n)
			} else {
				panic(fmt.Sprintf("Invalid text/Output input %q", i))
			}
		}
	}()

	return o
}

func (*output) Description() string {
	return "Outputs each packet, forwards it to the assigned channel."
}

func (o *output) Inputs() map[string]process.Input {
	return map[string]process.Input{
		"in": {Channel: o.in, IPType: process.AnyIP},
	}
}

func (*output) Outputs() map[string]process.Output {
	return map[string]process.Output{}
}
