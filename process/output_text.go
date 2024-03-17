package process

import "fmt"

type outputText struct {
	ins  map[string]Input
	outs map[string]Output
}

func OutputText(out chan<- string) Process {
	in := make(chan any, 1)

	go func() {
		for i := range in {
			if s, ok := i.(string); !ok {
				panic(fmt.Sprintf("Invalid input %q", i))
			} else {
				out <- s
			}
		}
	}()

	return &outputText{
		ins:  map[string]Input{"in": {Channel: in, IPType: StringIP}},
		outs: map[string]Output{},
	}
}

func (ot *outputText) Inputs() map[string]Input {
	return ot.ins
}

func (ot *outputText) Outputs() map[string]Output {
	return ot.outs
}
