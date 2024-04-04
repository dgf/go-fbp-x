package util_test

import (
	"strings"
	"testing"
	"time"

	"github.com/dgf/go-fbp-x/util"
)

func TestParseTimeISO8601(t *testing.T) {
	for _, tc := range []struct {
		name string
		in   any
		err  string
		exp  time.Duration
	}{
		{"zero reject", "", "zero", time.Duration(0)},
		{"invalid input", "fooS", "invalid", time.Duration(0)},
		{"invalid definition", false, "bool", time.Duration(0)},
		{"137 ms", "0.137S", "", time.Duration(137 * time.Millisecond)},
		{"1.37 ms", "1.37S", "", time.Duration(1*time.Second + 370*time.Millisecond)},
		{"17 seconds", "17S", "", time.Duration(17 * time.Second)},
		{"37 minutes", "37M", "", time.Duration(37 * time.Minute)},
		{"73 hours", "73H", "", time.Duration(73 * time.Hour)},
		{"all parts", "2H3M4S", "", time.Duration(2*time.Hour + 3*time.Minute + 4*time.Second)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			act, err := util.ParseTimeISO8601(tc.in)
			if len(tc.err) > 0 {
				if err == nil {
					t.Errorf("ParseTimeISO8601 no error, but should contain: %s", tc.err)
				} else if !strings.Contains(err.Error(), tc.err) {
					t.Errorf("ParseTimeISO8601 invalid error, got: %s, should contain: %s", err, tc.err)
				}
			} else if err != nil {
				t.Errorf("ParseTimeISO8601 unexpected error: %v", err)
			} else if act != tc.exp {
				t.Errorf("ParseTimeISO8601 unmatched duration, got: %s, want: %s", act, tc.exp)
			}
		})
	}
}
