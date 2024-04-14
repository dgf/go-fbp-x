package util_test

import (
	"strings"
	"testing"

	"github.com/dgf/go-fbp-x/pkg/util"
)

func TestCastAndUnescapeRaw(t *testing.T) {
	for _, tc := range []struct {
		name string
		sep  any
		exp  string
		err  string
	}{
		{"char dot error", '.', "", "value"},
		{"string slash", "/", "/", ""},
		{"string new line", "\n", "\n", ""},
		{"raw dash", `-`, "-", ""},
		{"raw new line", `\n`, "\n", ""},
		{"raw tab", `\t`, "\t", ""},
		{"raw carriage return", `\r`, "\r", ""},
		{"all raws multi", `\n\r\t\r\n`, "\n\r\t\r\n", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			act, err := util.CastAndUnescapeRaw(tc.sep)
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

func BenchmarkCastAndUnescapeRaw(b *testing.B) {
	for _, tc := range []struct {
		name string
		raw  string
	}{
		{"simple", `test`},
		{"breaks", `one\nwin\n\r\two`},
	} {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				util.CastAndUnescapeRaw(tc.raw)
			}
		})
	}
}
