package text_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/dgf/go-fbp-x/process/text"
)

func TestUnescape(t *testing.T) {
	for _, tc := range []struct {
		name string
		sep  any
		exp  string
		err  string
	}{
		{"char dot error", '.', "", "invalid"},
		{"string slash", "/", "/", ""},
		{"string new line", "\n", "\n", ""},
		{"raw dash", `-`, "-", ""},
		{"raw new line", `\n`, "\n", ""},
		{"raw tab", `\t`, "\t", ""},
		{"raw carriage return", `\r`, "\r", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			act, err := text.NewSplit().CastAndUnescape(tc.sep)
			if len(tc.err) > 0 {
				if err == nil {
					t.Errorf("Unescape error expected")
				} else if !strings.Contains(err.Error(), tc.err) {
					t.Errorf("Unescape error got: %q, should contain: %q", err, tc.err)
				}
			} else if err != nil {
				t.Errorf("Unescape error %q", err)
			}
			if act != tc.exp {
				t.Errorf("Unescape seperator got: %q, want: %q", act, tc.exp)
			}
		})
	}
}

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
			in := p.Inputs()["in"]
			sep := p.Inputs()["sep"]
			out := p.Outputs()["out"]
			act := make(chan []string, 1)

			go func() {
				sep.Channel <- tc.sep
				in.Channel <- tc.in
			}()

			go func() {
				a := make([]string, len(tc.out))
				for range len(tc.out) {
					s := <-out.Channel
					a = append(a, s.(string))
				}
				act <- a
			}()

			if reflect.DeepEqual(act, tc.out) {
				t.Errorf("Split got: %v, want: %v", act, tc.out)
			}
		})
	}
}
