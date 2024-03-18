package process

type counter struct {
	ins  map[string]Input
	outs map[string]Output
}

func Counter() Process {
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

func (sl *counter) Inputs() map[string]Input {
	return sl.ins
}

func (sl *counter) Outputs() map[string]Output {
	return sl.outs
}
