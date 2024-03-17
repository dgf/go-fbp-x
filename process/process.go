package process

type IPType int64

const (
	NumberIP IPType = iota
	StringIP
)

type Input struct {
	Stream chan<- any
	Kind   IPType
}

type Output struct {
	Stream <-chan any
	Kind   IPType
}

type Process interface {
	Inputs() map[string]Input
	Outputs() map[string]Output
}
