package core

import (
	"fmt"
	"math"
	"strconv"
	"time"
	"unicode"

	"github.com/dgf/go-fbp-x/process"
)

type tick struct {
	data chan any
	intv chan any
	out  chan any
}

func ParseTimeISO8601(t any) (time.Duration, error) {
	d := time.Duration(0)
	s, ok := t.(string)
	if !ok {
		return d, fmt.Errorf("invalid duration input %q", t)
	}

	n := ""
	for _, c := range s {
		switch c {
		case 'H':
			if i, err := strconv.Atoi(n); err != nil {
				return d, err
			} else {
				n = ""
				d += time.Duration(i * int(time.Hour))
			}
		case 'M':
			if i, err := strconv.Atoi(n); err != nil {
				return d, err
			} else {
				n = ""
				d += time.Duration(i * int(time.Minute))
			}
		case 'S':
			if f, err := strconv.ParseFloat(n, 32); err != nil {
				return d, err
			} else {
				n = ""
				sec, ms := math.Modf(f)
				d += time.Duration(int(sec)*int(time.Second) + int(math.Round(ms*1000))*int(time.Millisecond))
			}
		default:
			if unicode.IsNumber(c) || c == '.' {
				n += string(c)
				continue
			}
		}
	}

	if d == 0 {
		return d, fmt.Errorf("invalid definition %q causes a duration of zero", t)
	}
	return d, nil
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

		intv, err := ParseTimeISO8601(i)
		if err != nil {
			panic(err)
		}

		for {
			select {
			case i, ok := <-t.intv:
				if !ok {
					return
				}
				intv, err = ParseTimeISO8601(i)
				if err != nil {
					panic(err)
				}
			case data = <-t.data:
				if !ok {
					return
				}
			case <-time.After(intv):
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
