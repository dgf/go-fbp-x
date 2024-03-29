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
		{name: "output data input", out: []string{"test"}, fbp: "'test' -> IN Display(core/Output)"},
		{name: "unknown file", out: []string{"open not-found.txt: no such file or directory"}, fbp: `
                'not-found.txt' -> IN Read(fs/ReadFile) ERR -> IN Display(core/Output)`},
		{name: "slurp multi inputs", out: []string{"one", "two"}, fbp: `
                'one' -> IN Display(core/Output)
                'two' -> IN Display`},
		{name: "demux output", out: []string{"one", "one"}, fbp: `
                'one' -> IN Demux(core/Clone) OUT -> IN Display1(core/Output)
                Demux OUT -> IN Display2(core/Output)`},
		{name: "split string", out: []string{"one", "two"}, ord: true, fbp: `
                '|' -> SEP Split(text/Split)
                'one|two' -> IN Split OUT -> IN Display(core/Output)`},
		{name: "count lines of text file", out: []string{"1 one", "2 two", "3 three"}, ord: true, fbp: `
                '\n' -> SEP Split(text/Split)
                ' ' -> DATA Space(core/Kick)
                'testdata/three-lines.txt' -> IN Read(fs/ReadFile)
                Read OUT -> IN Split OUT -> IN Count(core/Count) OUT -> IN CountAndSpace(text/Append)
                Count OUT -> IN Space OUT -> AFFIX CountAndSpace OUT -> IN CountAndLine(text/Append)
                Split OUT -> AFFIX CountAndLine OUT -> IN Display(core/Output)
                Read ERR -> IN Display`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			out := make(chan string, 1)
			done := make(chan []string, 1)
			traces := make(chan network.Trace, 1)

			if graph, err := dsl.Parse(strings.NewReader(tc.fbp)); err != nil {
				t.Errorf("Parse failed: %v", err)
				return
			} else if err := network.NewNetwork(out).Run(graph, traces); err != nil {
				t.Errorf("Run failed: %v", err)
				return
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
