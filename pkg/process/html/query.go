package html

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/andybalholm/cascadia"
	"github.com/dgf/go-fbp-x/pkg/process"
	"golang.org/x/net/html"
)

type query struct {
	in  chan any
	sel chan any
	out chan any
	err chan any
}

func Query() process.Process {
	q := &query{
		in:  make(chan any, 1),
		sel: make(chan any, 1),
		out: make(chan any, 1),
		err: make(chan any, 1),
	}

	go func() {
		defer close(q.err)
		defer close(q.out)

		s, ok := <-q.sel
		if !ok {
			return
		}

		sel, ok := s.(string)
		if !ok {
			panic(fmt.Sprintf("Invalid html/Query selector %q", sel))
		}

		for {
			select {
			case s, ok := <-q.sel:
				if !ok {
					return
				}

				sel, ok = s.(string)
				if !ok {
					panic(fmt.Sprintf("Invalid html/Query selector %q", sel))
				}
			case i, ok := <-q.in:
				if !ok {
					return
				}

				in, ok := i.(string)
				if !ok {
					panic(fmt.Sprintf("Invalid html/Query input %v", i))
				}

				cSel, err := cascadia.Parse(sel)
				if err != nil {
					q.err <- fmt.Sprintf("parse of html/Query selector failed: %v", err)
					return
				}

				doc, err := html.Parse(strings.NewReader(in))
				if err != nil {
					q.err <- fmt.Sprintf("parse of html/Query document failed: %v", err)
					return
				}

				for _, n := range cascadia.QueryAll(doc, cSel) {
					var b bytes.Buffer
					if err := html.Render(&b, n); err != nil {
						q.err <- fmt.Sprintf("render html/Query output node failed: %v", err)
					} else {
						q.out <- b.String()
					}
				}
			}
		}
	}()

	return q
}

func (*query) Description() string {
	return "Outputs selected HTML elements."
}

func (q *query) Inputs() map[string]process.Input {
	return map[string]process.Input{
		"in":  {Channel: q.in, IPType: process.StringIP},
		"sel": {Channel: q.sel, IPType: process.StringIP},
	}
}

func (g *query) Outputs() map[string]process.Output {
	return map[string]process.Output{
		"err": {Channel: g.err, IPType: process.StringIP},
		"out": {Channel: g.out, IPType: process.StringIP},
	}
}
