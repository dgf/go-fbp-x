package network_test

import (
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/dgf/go-fbp-x/dsl"
	"github.com/dgf/go-fbp-x/network"
)

func TestRun(t *testing.T) {
	for _, tc := range []struct {
		name string
		fbp  string
		out  []string
		ord  bool
	}{
		{name: "output data input", out: []string{"test"}, fbp: "'test' -> IN Display(OutputText)"},
		{name: "unknown file", out: []string{"open not-found.txt: no such file or directory"}, fbp: `
                'not-found.txt' -> IN Read(ReadFile) ERROR -> IN Display(OutputText)`},
		{name: "slurp multi inputs", out: []string{"one", "two"}, fbp: `
                'one' -> IN Display(OutputText)
                'two' -> IN Display`},
		{name: "demux output", out: []string{"one", "one"}, fbp: `
                'one' -> IN Demux(Clone) OUT -> IN Display1(OutputText)
                Demux OUT -> IN Display2(OutputText)`},
		{name: "count lines of text file", out: []string{"1", "2", "3"}, ord: true, fbp: `
                'testdata/three-lines.txt' -> IN Read(ReadFile)
                Read OUT -> IN Split(SplitLines) OUT -> IN Count(Count)
                Count OUT -> IN Display(OutputText)
                Read ERROR -> IN Display`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			out := make(chan string, 1)
			done := make(chan []string, 1)
			traces := make(chan network.Trace, 1)

			if graph, err := dsl.Parse(strings.NewReader(tc.fbp)); err != nil {
				t.Errorf("Parse failed: %v", err)
			} else if network, err := network.Create(graph, out); err != nil {
				t.Errorf("Create failed: %v", err)
			} else {
				network.Run(traces)
			}

			go func() {
				for range traces {
					// discard
				}
			}()

			go func() {
				act := make([]string, len(tc.out))
				for i := range len(tc.out) {
					act[i] = <-out
				}
				done <- act
			}()

			select {
			case <-time.After(37 * time.Millisecond):
				t.Error("Timeout! Deadlock?")
			case act := <-done:
				if !tc.ord {
					slices.Sort(act)
					slices.Sort(tc.out)
				}
				if !reflect.DeepEqual(act, tc.out) {
					t.Errorf("Run failed got: %q, want: %q", act, tc.out)
				}
			}
		})
	}
}
