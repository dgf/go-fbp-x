package filesystem

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/dgf/go-fbp-x/process"
)

type readFile struct {
	in   chan any
	out  chan any
	errs chan any
}

func ReadFile() process.Process {
	rf := &readFile{
		in:   make(chan any, 1),
		out:  make(chan any, 1),
		errs: make(chan any, 1),
	}

	go func() {
		for i := range rf.in {
			if s, ok := i.(string); !ok {
				panic(fmt.Sprintf("Invalid fs/ReadFile input %q", i))
			} else if file, err := os.Open(s); err != nil {
				rf.errs <- err.Error()
			} else if data, err := io.ReadAll(file); err != nil {
				rf.errs <- err.Error()
			} else {
				rf.out <- strings.TrimRight(strings.TrimRight(string(data), "\n"), "\r")
			}
		}
	}()

	return rf
}

func (*readFile) Description() string {
	return "Reads a file and outputs the content as string."
}

func (rf *readFile) Inputs() map[string]process.Input {
	return map[string]process.Input{
		"in": {Channel: rf.in, IPType: process.StringIP},
	}
}

func (rf *readFile) Outputs() map[string]process.Output {
	return map[string]process.Output{
		"out": {Channel: rf.out, IPType: process.StringIP},
		"err": {Channel: rf.errs, IPType: process.StringIP},
	}
}
