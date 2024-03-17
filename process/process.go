package process

type IPType int64

const (
	NumberIP IPType = iota
	StringIP
	StringSliceIP
	AnySliceIP
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
