package process

import (
	"fmt"
	"io"
	"os"
)

type readFile struct {
	ins  map[string]chan<- string
	outs map[string]<-chan string
}

func ReadFile() Process {
	in := make(chan string, 1)
	out := make(chan string, 1)
	errs := make(chan string, 1)

	go func() {
		for i := range in {
			fmt.Println("read file:", i)
			if file, err := os.Open(i); err != nil {
				errs <- err.Error()
			} else if data, err := io.ReadAll(file); err != nil {
				errs <- err.Error()
			} else {
				out <- string(data)
			}
		}
	}()

	return &readFile{
		ins:  map[string]chan<- string{"in": in},
		outs: map[string]<-chan string{"out": out, "error": errs},
	}
}

func (rf *readFile) Inputs() map[string]chan<- string {
	return rf.ins
}

func (rf *readFile) Outputs() map[string]<-chan string {
	return rf.outs
}
