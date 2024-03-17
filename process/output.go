package process

type output struct {
	ins  map[string]chan<- string
	outs map[string]<-chan string
}

func Output(out chan<- string) Process {
	in := make(chan string, 1)

	go func() {
		for i := range in {
			out <- i
		}
	}()

	return &output{
		ins:  map[string]chan<- string{"in": in},
		outs: map[string]<-chan string{},
	}
}

func (o *output) Inputs() map[string]chan<- string {
	return o.ins
}

func (o *output) Outputs() map[string]<-chan string {
	return o.outs
}
