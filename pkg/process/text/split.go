package text

import (
	"fmt"
	"strings"

	"github.com/dgf/go-fbp-x/pkg/process"
	"github.com/dgf/go-fbp-x/pkg/util"
)

type split struct {
	in  chan any
	out chan any
	sep chan any
}

func Split() process.Process {
	s := &split{
		in:  make(chan any, 1),
		out: make(chan any, 1),
		sep: make(chan any, 1),
	}

	go func() {
		defer close(s.out)

		if err := s.Process(); err != nil {
			panic(fmt.Sprintf("Invalid core/Split %v", err))
		}
	}()

	return s
}

func (s *split) Process() error {
	if firstSep, ok := <-s.sep; !ok {
		return nil
	} else if currentSep, err := util.CastAndUnescapeRaw(firstSep); err != nil { // requires a seperator upfront
		return fmt.Errorf("separator: %w", err)
	} else {
		for {
			select {
			case nextSep, ok := <-s.sep: // replace seperator
				if !ok {
					return nil
				}
				currentSep, err = util.CastAndUnescapeRaw(nextSep)
				if err != nil {
					return fmt.Errorf("separator: %w", err)
				}
			case in, ok := <-s.in:
				if !ok {
					return nil
				}
				if is, ok := in.(string); !ok {
					return fmt.Errorf("input %q", in)
				} else {
					for _, part := range strings.Split(is, currentSep) {
						s.out <- part
					}
				}
			}
		}
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
