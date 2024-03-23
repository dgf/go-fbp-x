package process

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

func IsCompatibleIPType(source IPType, target IPType) bool {
	return source == target || target == AnyIP
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
