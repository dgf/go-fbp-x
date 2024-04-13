package util

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func ParseTimeISO8601(t any) (time.Duration, error) {
	d := time.Duration(0)
	s, ok := t.(string)
	if !ok {
		return d, fmt.Errorf("invalid duration input %q", t)
	}

	n := strings.Builder{}
	for _, c := range s {
		switch c {
		case 'H':
			if i, err := strconv.Atoi(n.String()); err != nil {
				return d, err
			} else {
				n.Reset()
				d += time.Duration(int64(i) * int64(time.Hour))
			}
		case 'M':
			if i, err := strconv.Atoi(n.String()); err != nil {
				return d, err
			} else {
				n.Reset()
				d += time.Duration(int64(i) * int64(time.Minute))
			}
		case 'S':
			if f, err := strconv.ParseFloat(n.String(), 32); err != nil {
				return d, err
			} else {
				n.Reset()
				sec, ms := math.Modf(f)
				d += time.Duration(int64(sec)*int64(time.Second) + int64(math.Round(ms*1000))*int64(time.Millisecond))
			}
		default:
			if unicode.IsNumber(c) || c == '.' {
				n.WriteRune(c)
				continue
			}
		}
	}

	if d == 0 {
		return d, fmt.Errorf("invalid definition %q causes a duration of zero", t)
	}
	return d, nil
}
