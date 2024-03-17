package network_test

import (
	"strings"
	"testing"
	"time"

	"github.com/dgf/go-fbp-x/dsl"
	"github.com/dgf/go-fbp-x/network"
)

func TestRun(t *testing.T) {
	for _, tc := range []struct {
		name string
		exp  string
		fbp  string
	}{
		{"output data input", "test", "'test' -> IN Display(OutputText)"},
		{"unknown file", "open not-found.txt: no such file or directory", `
                'not-found.txt' -> IN Read(ReadFile) ERROR -> IN Display(OutputText)`},
		{"count lines of text file", "3", `
                'testdata/three-lines.txt' -> IN Read(ReadFile)
                Read OUT -> IN Split(SplitLines) OUT -> IN Count(Counter)
                Count OUT -> IN Display(OutputText)
                Read ERROR -> IN Display`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			out := make(chan string, 1)

			if graph, err := dsl.Parse(strings.NewReader(tc.fbp)); err != nil {
				t.Errorf("Parse failed: %v", err)
			} else if network, err := network.Create(graph, out); err != nil {
				t.Errorf("Create failed: %v", err)
			} else if err := network.Run(); err != nil {
				t.Errorf("Run failed: %v", err)
			}

			select {
			case <-time.After(37 * time.Millisecond):
				t.Error("Timeout Run")
			case act := <-out:
				if act != tc.exp {
					t.Errorf("Run failed got: %q, want: %q", act, tc.exp)
				}
			}
		})
	}
}
