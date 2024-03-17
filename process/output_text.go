package process

import (
	"fmt"
	"strconv"
)

type outputText struct {
	ins  map[string]Input
	outs map[string]Output
}

func OutputText(out chan<- string) Process {
	in := make(chan any, 1)

	go func() {
		for i := range in {
			if s, ok := i.(string); ok {
				out <- s
			} else if n, ok := i.(int); ok {
				out <- strconv.Itoa(n)
			} else {
				panic(fmt.Sprintf("Invalid input %q", i))
			}
		}
	}()

	return &outputText{
		ins:  map[string]Input{"in": {Channel: in, IPType: AnyIP}},
		outs: map[string]Output{},
	}
}

func (ot *outputText) Inputs() map[string]Input {
	return ot.ins
}

func (ot *outputText) Outputs() map[string]Output {
	return ot.outs
}
