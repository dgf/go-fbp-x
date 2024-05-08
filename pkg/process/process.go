package process

import "errors"

type IPType int64

const (
	AnyIP IPType = iota
	NumberIP
	StringIP
)

type Input struct {
	Channel chan<- any
	IPType
}

type Output struct {
	Channel <-chan any
	IPType
}

type Process interface {
	Description() string
	Inputs() map[string]Input
	Outputs() map[string]Output
}

func ConvertIP(source, target IPType) (func(any) any, error) {
	if source == target || target == AnyIP {
		return func(in any) any {
			return in
		}, nil
	}
	if source == AnyIP && target == StringIP {
		return func(in any) any {
			return in.(string)
		}, nil
	}
	return nil, errors.New("no converter available")
}

func (ipt IPType) String() string {
	switch ipt {
	case AnyIP:
		return "any"
	case NumberIP:
		return "number"
	case StringIP:
		return "string"
	}
	return "unknown"
}
