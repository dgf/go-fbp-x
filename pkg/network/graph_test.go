package network_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/dgf/go-fbp-x/pkg/network"
)

func ExampleString() {
	g := network.Graph{
		Components: map[string]string{
			"one":   "DoSome",
			"two":   "DoOther",
			"three": "DoElse",
		},
		Connections: []network.Connection{
			{Data: "data", Target: network.Link{Port: "in", Component: "one"}},
			{Source: network.Link{Port: "out", Component: "one"}, Target: network.Link{Port: "in", Component: "two"}},
			{Source: network.Link{Port: "out", Component: "one"}, Target: network.Link{Port: "in", Index: 1, Component: "three"}},
			{Source: network.Link{Port: "out", Index: 1, Component: "two"}, Target: network.Link{Port: "in", Index: 2, Component: "three"}},
		},
	}

	fmt.Println(g.String())
	// Output:
	//
	// components:
	// one > DoSome
	// three > DoElse
	// two > DoOther
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
		conns []network.Connection
	}{
		{"empty", []string{}, []network.Connection{}},
		{"data packet with one target", []string{"foo"}, []network.Connection{
			{Data: "data", Target: network.Link{Component: "foo"}},
		}},
		{"one target as source", []string{}, []network.Connection{
			{Source: network.Link{Component: "foo"}, Target: network.Link{Component: "foo"}},
		}},
		{"one target and one source", []string{"bar"}, []network.Connection{
			{Source: network.Link{Component: "foo"}, Target: network.Link{Component: "bar"}},
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			act := network.Leaves(tc.conns)
			if !reflect.DeepEqual(act, tc.exits) {
				t.Errorf("Wrong exit detection got: %v, want: %v", act, tc.exits)
			}
		})
	}
}
