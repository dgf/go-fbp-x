package dsl_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/dgf/go-fbp-x/pkg/dsl"
)

func ExampleString() {
	g := dsl.Graph{
		Components: map[string]dsl.Process{
			"one":   {Name: "DoSome"},
			"two":   {Name: "DoOther", Meta: map[string]string{"key": "value", "num": "123", "flag": "true"}},
			"three": {Name: "DoElse"},
		},
		Connections: []dsl.Connection{
			{Data: "data", Target: dsl.Link{Port: "in", Component: "one"}},
			{Source: dsl.Link{Port: "out", Component: "one"}, Target: dsl.Link{Port: "in", Component: "two"}},
			{Source: dsl.Link{Port: "out", Component: "one"}, Target: dsl.Link{Port: "in", Index: 1, Component: "three"}},
			{Source: dsl.Link{Port: "out", Index: 1, Component: "two"}, Target: dsl.Link{Port: "in", Index: 2, Component: "three"}},
		},
	}

	fmt.Println(g.String())
	// Output:
	//
	// components:
	// one > DoSome
	// three > DoElse
	// two > DoOther:flag=true,key=value,num=123
	//
	// connections:
	// data > in one
	// one out > in two
	// one out > in[1] three
	// two out[1] > in[2] three
}

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
