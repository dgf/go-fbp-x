package network_test

import (
	"reflect"
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
	}{
		{name: "output data input", out: []string{"test"}, fbp: "'test' -> IN Display(OutputText)"},
		{name: "unknown file", out: []string{"open not-found.txt: no such file or directory"}, fbp: `
                'not-found.txt' -> IN Read(ReadFile) ERROR -> IN Display(OutputText)`},
		{name: "count lines of text file", out: []string{"1", "2", "3"}, fbp: `
                'testdata/three-lines.txt' -> IN Read(ReadFile)
                Read OUT -> IN Split(SplitLines) OUT -> IN Count(Counter)
                Count OUT -> IN Display(OutputText)
                Read ERROR -> IN Display`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			out := make(chan string, 1)
			act := make(chan []string, 1)

			if graph, err := dsl.Parse(strings.NewReader(tc.fbp)); err != nil {
				t.Errorf("Parse failed: %v", err)
			} else if network, err := network.Create(graph, out); err != nil {
				t.Errorf("Create failed: %v", err)
			} else {
				network.Run()
			}

			go func() {
				a := make([]string, len(tc.out))
				for i := range len(tc.out) {
					a[i] = <-out
				}
				act <- a
			}()

			select {
			case <-time.After(37 * time.Millisecond):
				t.Error("Timeout Run")
			case a := <-act:
				if !reflect.DeepEqual(a, tc.out) {
					t.Errorf("Run failed got: %q, want: %q", a, tc.out)
				}
			}
		})
	}
}
