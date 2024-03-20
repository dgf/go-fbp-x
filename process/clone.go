package process

type clone struct {
	ins  map[string]Input
	outs map[string]Output
}

func Clone() Process {
	in := make(chan any, 1)
	out := make(chan any, 1)

	go func() {
		for i := range in {
			out <- i
		}
	}()

	return &outputText{
		ins:  map[string]Input{"in": {Channel: in, IPType: AnyIP}},
		outs: map[string]Output{"out": {Channel: out, IPType: AnyIP}},
	}
}

func (c *clone) Inputs() map[string]Input {
	return c.ins
}

func (c *clone) Outputs() map[string]Output {
	return c.outs
}
