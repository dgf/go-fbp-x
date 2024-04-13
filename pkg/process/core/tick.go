package core

import (
	"fmt"
	"time"

	"github.com/dgf/go-fbp-x/pkg/process"
	"github.com/dgf/go-fbp-x/pkg/util"
)

type tick struct {
	data chan any
	intv chan any
	out  chan any
}

func Tick() process.Process {
	t := &tick{
		data: make(chan any, 1),
		intv: make(chan any, 1),
		out:  make(chan any, 1),
	}

	go func() {
		defer close(t.out)

		data, ok := <-t.data
		if !ok {
			return
		}

		i, ok := <-t.intv
		if !ok {
			return
		}

		intv, err := util.ParseTimeISO8601(i)
		if err != nil {
			panic(fmt.Sprintf("Invalid core/Tick interval: %v", err))
		}

		timer := time.NewTimer(intv)
		for {
			timer.Reset(intv)

			select {
			case i, ok := <-t.intv:
				if !ok {
					return
				}
				intv, err = util.ParseTimeISO8601(i)
				if err != nil {
					panic(fmt.Sprintf("Invalid core/Tick interval: %v", err))
				}
			case data, ok = <-t.data:
				if !ok {
					return
				}
			case <-timer.C:
				t.out <- data
			}
		}
	}()

	return t
}

func (*tick) Description() string {
	return "Sends the data packet on every interval tick."
}

func (t *tick) Inputs() map[string]process.Input {
	return map[string]process.Input{
		"data": {Channel: t.data, IPType: process.AnyIP},
		"intv": {Channel: t.intv, IPType: process.StringIP},
	}
}

func (t *tick) Outputs() map[string]process.Output {
	return map[string]process.Output{
		"out": {Channel: t.out, IPType: process.AnyIP},
	}
}
