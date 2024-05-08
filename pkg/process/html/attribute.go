package html

import (
	"fmt"
	"strings"

	"github.com/dgf/go-fbp-x/pkg/process"
	"golang.org/x/net/html"
)

type attribute struct {
	in  chan any
	sel chan any
	out chan any
	err chan any
}

func attr(in, sel string) (string, error) {
	t := html.NewTokenizer(strings.NewReader(in))

	tt := t.Next()
	if tt == html.ErrorToken {
		return "", fmt.Errorf("attribute %q not found in %v", sel, in)
	}

	for _, a := range t.Token().Attr {
		if a.Key == sel {
			return a.Val, nil
		}
	}

	return "", fmt.Errorf("attribute %q not found in %v", sel, in)
}

func Attribute() process.Process {
	a := &attribute{
		in:  make(chan any, 1),
		sel: make(chan any, 1),
		out: make(chan any, 1),
		err: make(chan any, 1),
	}

	go func() {
		defer close(a.err)
		defer close(a.out)

		s, ok := <-a.sel
		if !ok {
			return
		}

		sel, ok := s.(string)
		if !ok {
			panic(fmt.Sprintf("Invalid html/Attribute selector %q", sel))
		}

		for {
			select {
			case s, ok := <-a.sel:
				if !ok {
					return
				}

				sel, ok = s.(string)
				if !ok {
					panic(fmt.Sprintf("Invalid html/Attribute selector %q", sel))
				}
			case i, ok := <-a.in:
				if !ok {
					return
				}

				in, ok := i.(string)
				if !ok {
					panic(fmt.Sprintf("Invalid html/Attribute input %v", i))
				}

				if v, err := attr(in, sel); err != nil {
					a.err <- err.Error()
				} else {
					a.out <- v
				}
			}
		}
	}()

	return a
}

func (*attribute) Description() string {
	return "Outputs selected HTML attribute."
}

func (a *attribute) Inputs() map[string]process.Input {
	return map[string]process.Input{
		"in":  {Channel: a.in, IPType: process.StringIP},
		"sel": {Channel: a.sel, IPType: process.StringIP},
	}
}

func (a *attribute) Outputs() map[string]process.Output {
	return map[string]process.Output{
		"err": {Channel: a.err, IPType: process.StringIP},
		"out": {Channel: a.out, IPType: process.StringIP},
	}
}
