package network_test

import (
	"context"
	"reflect"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/dgf/go-fbp-x/pkg/dsl"
	"github.com/dgf/go-fbp-x/pkg/network"
	"github.com/dgf/go-fbp-x/pkg/process"
	"github.com/dgf/go-fbp-x/pkg/process/core"
	"github.com/dgf/go-fbp-x/pkg/process/filesystem"
	"github.com/dgf/go-fbp-x/pkg/process/text"
)

func withoutMeta(fn func() process.Process) func(map[string]string) process.Process {
	return func(map[string]string) process.Process {
		return fn()
	}
}

func newFactory(out chan<- string) process.Factory {
	return process.NewFactory(map[string]func(map[string]string) process.Process{
		"core/Clone":  withoutMeta(core.Clone),
		"core/Count":  withoutMeta(core.Count),
		"core/Kick":   withoutMeta(core.Kick),
		"core/Output": func(map[string]string) process.Process { return core.Output(out) },
		"core/Tick":   withoutMeta(core.Tick),
		"fs/ReadFile": withoutMeta(filesystem.ReadFile),
		"text/Append": text.Append,
		"text/Split":  withoutMeta(text.Split),
	})
}

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
		{name: "count two ticks", out: []string{"1bang", "2bang"}, ord: true, fbp: `
                '0.001S' -> INTV Ticker(core/Tick)
                'bang' -> DATA Ticker OUT -> IN Count(core/Count) OUT -> IN Append(text/Append)
                Ticker OUT -> AFFIX Append OUT -> IN Display(core/Output)`},
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
			graph, err := dsl.Parse(strings.NewReader(tc.fbp))
			if err != nil {
				t.Errorf("Parse failed: %v", err)
				return
			}

			wg := sync.WaitGroup{}
			out := make(chan string, 1)
			done := make(chan []string, 1)
			ctx, cancelRun := context.WithCancel(context.Background())

			traces := make(chan network.Trace, 1)
			defer close(traces)
			go func() {
				for range traces {
					// discard
				}
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()

				if err := network.NewNetwork(newFactory(out)).Run(ctx, graph, traces); err != nil {
					t.Errorf("Run failed: %v", err)
					return
				}
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()

				act := make([]string, len(tc.out))
				for i := range len(tc.out) {
					act[i] = <-out
				}
				cancelRun()
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

			wg.Wait()
		})
	}
}
