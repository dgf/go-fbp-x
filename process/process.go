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
	Inputs() map[string]Input
	Outputs() map[string]Output
}

func IsCompatibleIPType(source IPType, target IPType) bool {
	return source == target || target == AnyIP
}
