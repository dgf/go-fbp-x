package dsl_test

import (
	"reflect"
	"testing"

	"github.com/dgf/go-fbp-x/dsl"
)

func TestLeaves(t *testing.T) {
	for _, tc := range []struct {
		name  string
		exits []string
		conns []dsl.Connection
	}{
		{"empty", []string{}, []dsl.Connection{}},
		{"data packet with one target", []string{"foo"}, []dsl.Connection{
			{Data: "data", Target: dsl.Link{Component: "foo"}},
		}},
		{"one target as source", []string{}, []dsl.Connection{
			{Source: dsl.Link{Component: "foo"}, Target: dsl.Link{Component: "foo"}},
		}},
		{"one target and one source", []string{"bar"}, []dsl.Connection{
			{Source: dsl.Link{Component: "foo"}, Target: dsl.Link{Component: "bar"}},
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			act := dsl.Leaves(tc.conns)
			if !reflect.DeepEqual(act, tc.exits) {
				t.Errorf("Wrong exit detection got: %v, want: %v", act, tc.exits)
			}
		})
	}
}
