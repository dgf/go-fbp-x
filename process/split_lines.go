package process

import (
	"fmt"
	"strings"
)

type splitLines struct {
	ins  map[string]Input
	outs map[string]Output
}

func SplitLines() Process {
	in := make(chan any, 1)
	out := make(chan any, 1)

	go func() {
		for i := range in {
			if s, ok := i.(string); !ok {
				panic(fmt.Sprintf("Invalid input %q", i))
			} else {
				out <- strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n")
			}
		}
	}()

	return &outputText{
		ins:  map[string]Input{"in": {Channel: in, IPType: StringIP}},
		outs: map[string]Output{"out": {Channel: out, IPType: StringSliceIP}},
	}
}

func (sl *splitLines) Inputs() map[string]Input {
	return sl.ins
}

func (sl *splitLines) Outputs() map[string]Output {
	return sl.outs
}
