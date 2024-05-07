package http

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dgf/go-fbp-x/pkg/process"
	"github.com/dgf/go-fbp-x/pkg/util"
)

type get struct {
	err     chan any
	out     chan any
	timeout chan any
	url     chan any
}

func fetch(url string, timeout time.Duration) (string, error) {
	client := http.Client{Timeout: timeout}

	if req, err := http.NewRequest(http.MethodGet, url, nil); err != nil {
		return "", fmt.Errorf("could not create http/Get request: %w", err)
	} else if res, err := client.Do(req); err != nil {
		return "", fmt.Errorf("could not do http/Get request: %w", err)
	} else if body, err := io.ReadAll(res.Body); err != nil {
		return "", fmt.Errorf("could not read http/Get response: %w", err)
	} else {
		defer res.Body.Close()
		return string(body), nil
	}
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

				if body, err := fetch(url, timeout); err != nil {
					g.err <- err.Error()
				} else {
					g.out <- string(body)
				}
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
