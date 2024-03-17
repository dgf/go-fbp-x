package process

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type readFile struct {
	ins  map[string]Input
	outs map[string]Output
}

func ReadFile() Process {
	in := make(chan any, 1)
	out := make(chan any, 1)
	errs := make(chan any, 1)

	go func() {
		for i := range in {
			if s, ok := i.(string); !ok {
				panic(fmt.Sprintf("Invalid input %q", i))
			} else if file, err := os.Open(s); err != nil {
				errs <- err.Error()
			} else if data, err := io.ReadAll(file); err != nil {
				errs <- err.Error()
			} else {
				out <- strings.TrimRight(strings.TrimRight(string(data), "\n"), "\r")
			}
		}
	}()

	return &readFile{
		ins: map[string]Input{
			"in": {Channel: in, IPType: StringIP},
		},
		outs: map[string]Output{
			"out":   {Channel: out, IPType: StringIP},
			"error": {Channel: errs, IPType: StringIP},
		},
	}
}

func (rf *readFile) Inputs() map[string]Input {
	return rf.ins
}

func (rf *readFile) Outputs() map[string]Output {
	return rf.outs
}
