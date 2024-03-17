package process_test

import (
	"strings"
	"testing"
	"time"

	"github.com/dgf/go-fbp-x/dsl"
	"github.com/dgf/go-fbp-x/process"
)

func TestRun(t *testing.T) {
	for _, tc := range []struct {
		name string
		fbp  string
		exp  string
	}{
		{"output data input", "'test' -> IN Display(Output)", "test"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			out := make(chan string, 1)

			if graph, err := dsl.Parse(strings.NewReader(tc.fbp)); err != nil {
				t.Errorf("Parse failed: %v", err)
			} else if network, err := process.Create(graph, out); err != nil {
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
