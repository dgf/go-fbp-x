package html

import (
	"fmt"
	"strings"

	"github.com/dgf/go-fbp-x/pkg/process"
	"golang.org/x/net/html"
)

type text struct {
	in  chan any
	out chan any
}

func innerText(in string) string {
	t := html.NewTokenizer(strings.NewReader(in))
	for {
		tt := t.Next()
		if tt == html.ErrorToken {
			return ""
		}
		if tt == html.TextToken {
			return strings.TrimSpace(t.Token().Data)
		}
	}
}

func Text() process.Process {
	t := &text{
		in:  make(chan any, 1),
		out: make(chan any, 1),
	}

	go func() {
		defer close(t.out)

		for i := range t.in {
			if in, ok := i.(string); !ok {
				panic(fmt.Sprintf("Invalid html/Attribute input %v", i))
			} else {
				t.out <- innerText(in)
			}
		}
	}()

	return t
}

func (*text) Description() string {
	return "Outputs selected HTML text."
}

func (t *text) Inputs() map[string]process.Input {
	return map[string]process.Input{
		"in": {Channel: t.in, IPType: process.StringIP},
	}
}

func (t *text) Outputs() map[string]process.Output {
	return map[string]process.Output{
		"out": {Channel: t.out, IPType: process.StringIP},
	}
}
