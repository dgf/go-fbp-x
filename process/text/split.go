package text

import (
	"fmt"
	"strings"

	"github.com/dgf/go-fbp-x/process"
)

type split struct {
	in       chan any
	out      chan any
	sep      chan any
	replacer *strings.Replacer
}

func NewSplit() *split {
	return &split{
		in:       make(chan any, 1),
		out:      make(chan any, 1),
		sep:      make(chan any, 1),
		replacer: strings.NewReplacer(`\n`, "\n", `\r`, "\r", `\t`, "\t"),
	}
}

const splitErrPrefix = "Invalid core/Split"

func Split() process.Process {
	s := NewSplit()

	go func() {
		if currentSep, err := s.CastAndUnescape(<-s.sep); err != nil { // requires a seperator upfront
			panic(fmt.Sprintf("%s sep input: %v", splitErrPrefix, err))
		} else {
			for {
				select {
				case lastSep := <-s.sep: // replace seperator
					currentSep, err = s.CastAndUnescape(lastSep)
					if err != nil {
						panic(fmt.Sprintf("%s sep input: %v", splitErrPrefix, err))
					}
				case i := <-s.in:
					if is, ok := i.(string); !ok {
						panic(fmt.Sprintf("%s input %q", splitErrPrefix, i))
					} else {
						for _, l := range strings.Split(is, currentSep) {
							s.out <- l
						}
					}
				}
			}
		}
	}()

	return s
}

func (s *split) CastAndUnescape(sep any) (string, error) {
	if currentSep, ok := sep.(string); !ok {
		return "", fmt.Errorf("invalid seperator %q", sep)
	} else {
		return s.replacer.Replace(currentSep), nil
	}
}

func (*split) Description() string {
	return "Splits inputs by seperator."
}

func (s *split) Inputs() map[string]process.Input {
	return map[string]process.Input{
		"in":  {Channel: s.in, IPType: process.StringIP},
		"sep": {Channel: s.sep, IPType: process.StringIP},
	}
}

func (s *split) Outputs() map[string]process.Output {
	return map[string]process.Output{
		"out": {Channel: s.out, IPType: process.StringIP},
	}
}
