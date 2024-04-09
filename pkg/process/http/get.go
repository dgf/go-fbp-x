package http

import (
	"fmt"
	"io"
	"net/http"

	"github.com/dgf/go-fbp-x/pkg/process"
	"github.com/dgf/go-fbp-x/pkg/util"
)

type get struct {
	err     chan any
	out     chan any
	timeout chan any
	url     chan any
}

func Get() process.Process {
	g := &get{
		err:     make(chan any, 1),
		out:     make(chan any, 1),
		timeout: make(chan any, 1),
		url:     make(chan any, 1),
	}

	go func() {
		defer close(g.err)
		defer close(g.out)

		t, ok := <-g.timeout
		if !ok {
			return
		}

		timeout, err := util.ParseTimeISO8601(t)
		if err != nil {
			panic(fmt.Sprintf("Invalid http/Get timeout: %v", err))
		}

		for {
			select {
			case t, ok := <-g.timeout:
				if !ok {
					return
				}

				timeout, err = util.ParseTimeISO8601(t)
				if err != nil {
					panic(fmt.Sprintf("Invalid http/Get timeout: %v", err))
				}
			case u, ok := <-g.url:
				if !ok {
					return
				}

				url, ok := u.(string)
				if !ok {
					panic(fmt.Sprintf("Invalid http/Get url %v", u))
				}

				req, err := http.NewRequest(http.MethodGet, url, nil)
				if err != nil {
					g.err <- fmt.Sprintf("could not create http/Get request: %v", err)
					return
				}

				client := http.Client{
					Timeout: timeout,
				}

				resp, err := client.Do(req)
				if err != nil {
					g.err <- fmt.Sprintf("could not do http/Get request: %v", err)
					return
				}

				body, err := io.ReadAll(resp.Body)
				if err != nil {
					g.err <- fmt.Sprintf("could not read http/Get response: %v", err)
					return
				}

				g.out <- string(body)
			}
		}
	}()

	return g
}

func (*get) Description() string {
	return "Outputs the result of an HTTP GET call."
}

func (g *get) Inputs() map[string]process.Input {
	return map[string]process.Input{
		"url":     {Channel: g.url, IPType: process.StringIP},
		"timeout": {Channel: g.timeout, IPType: process.StringIP},
	}
}

func (g *get) Outputs() map[string]process.Output {
	return map[string]process.Output{
		"err": {Channel: g.err, IPType: process.StringIP},
		"out": {Channel: g.out, IPType: process.StringIP},
	}
}
