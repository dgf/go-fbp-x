package process

import (
	"reflect"
)

type counter struct {
	ins  map[string]Input
	outs map[string]Output
}

func Counter() Process {
	in := make(chan any, 1)
	out := make(chan any, 1)

	go func() {
		for i := range in {
			out <- reflect.ValueOf(i).Len()
		}
	}()

	return &outputText{
		ins:  map[string]Input{"in": {Channel: in, IPType: AnySliceIP}},
		outs: map[string]Output{"out": {Channel: out, IPType: NumberIP}},
	}
}

func (sl *counter) Inputs() map[string]Input {
	return sl.ins
}

func (sl *counter) Outputs() map[string]Output {
	return sl.outs
}
