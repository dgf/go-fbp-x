package text_test

import (
	"reflect"
	"testing"

	"github.com/dgf/go-fbp-x/pkg/process/text"
)

func TestSplit(t *testing.T) {
	for _, tc := range []struct {
		name string
		sep  any
		in   string
		out  []string
	}{
		{"split by string slash", "/", "foo/bar", []string{"foo", "bar"}},
		{"split by raw dot", `.`, "one.two", []string{"one", "two"}},
		{"split by raw line break", `\n`, "line1\nline2", []string{"line1", "line2"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			p := text.Split()
			out := p.Outputs()["out"]

			in := p.Inputs()["in"]
			defer close(in.Channel)
			sep := p.Inputs()["sep"]
			defer close(sep.Channel)

			done := make(chan []string, 1)
			defer close(done)

			go func() {
				sep.Channel <- tc.sep
				in.Channel <- tc.in
			}()

			go func() {
				act := make([]string, len(tc.out))
				for i := range len(tc.out) {
					o := <-out.Channel
					act[i] = o.(string)
				}
				done <- act
			}()

			act := <-done
			if !reflect.DeepEqual(act, tc.out) {
				t.Errorf("Split got: %v, want: %v", act, tc.out)
			}
		})
	}
}

func BenchmarkSplit(b *testing.B) {
	p := text.Split()

	in := p.Inputs()["in"]
	defer close(in.Channel)
	sep := p.Inputs()["sep"]
	defer close(sep.Channel)

	sep.Channel <- `test`

	for _, tc := range []struct {
		name string
		sep  string
	}{
		{"test", `test`},
		{"line breaks", `\n\r\t`},
	} {
		b.Run("separator "+tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				sep.Channel <- tc.sep
			}
		})
	}

	out := p.Outputs()["out"]
	go func() {
		for {
			<-out.Channel
		}
	}()

	for _, tc := range []struct {
		name string
		sep  string
		in   string
	}{
		{"line breaks", `\n`, "one\ntwo\nthree"},
		{"multi dots", `.`, "1.2.3.4.5.6.7.8.9"},
	} {
		sep.Channel <- tc.sep
		b.Run("input "+tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				in.Channel <- tc.in
			}
		})
	}
}
