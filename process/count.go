package process

type count struct {
	ins  map[string]Input
	outs map[string]Output
}

func Count() Process {
	in := make(chan any, 1)
	out := make(chan any, 1)
	count := 0

	go func() {
		for range in {
			count++
			out <- count
		}
	}()

	return &outputText{
		ins:  map[string]Input{"in": {Channel: in, IPType: AnyIP}},
		outs: map[string]Output{"out": {Channel: out, IPType: NumberIP}},
	}
}

func (c *count) Inputs() map[string]Input {
	return c.ins
}

func (c *count) Outputs() map[string]Output {
	return c.outs
}
