package text

import (
	"fmt"

	"github.com/dgf/go-fbp-x/process"
)

type append struct {
	in    chan any
	affix chan any
	out   chan any
}

func Append() process.Process {
	a := &append{
		in:    make(chan any, 1),
		affix: make(chan any, 1),
		out:   make(chan any, 1),
	}

	go func() {
		defer close(a.out)

		for {
			in, ok := <-a.in
			if !ok {
				return
			}

			affix, ok := <-a.affix
			if !ok {
				return
			}

			a.out <- fmt.Sprintf("%v%v", in, affix)
		}
	}()

	return a
}

func (*append) Description() string {
	return "Appends affix to input."
}

func (a *append) Inputs() map[string]process.Input {
	return map[string]process.Input{
		"in":    {Channel: a.in, IPType: process.AnyIP},
		"affix": {Channel: a.affix, IPType: process.AnyIP},
	}
}

func (a *append) Outputs() map[string]process.Output {
	return map[string]process.Output{
		"out": {Channel: a.out, IPType: process.StringIP},
	}
}
