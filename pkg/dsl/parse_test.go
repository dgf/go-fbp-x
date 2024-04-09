package dsl_test

import (
	"os"
	"strings"
	"testing"

	"github.com/dgf/go-fbp-x/pkg/dsl"
	"github.com/dgf/go-fbp-x/pkg/network"
)

func DataTarget(data, port, component string) network.Connection {
	return network.Connection{Data: data, Target: network.Link{Port: port, Component: component}}
}

func SourceTargetIndex(sComp, sPort string, sIndex, tIndex int, tPort, tComp string) network.Connection {
	return network.Connection{
		Source: network.Link{Port: sPort, Index: sIndex, Component: sComp},
		Target: network.Link{Port: tPort, Index: tIndex, Component: tComp},
	}
}

func SourceTarget(sComp, sPort, tPort, tComp string) network.Connection {
	return SourceTargetIndex(sComp, sPort, 0, 0, tPort, tComp)
}

func TestParse(t *testing.T) {
	for _, tc := range []struct {
		name  string
		fbp   string
		graph network.Graph
	}{
		{"empty", "", network.Graph{}},
		{"one comment", "# comment", network.Graph{}},
		{"one component", "in call", network.Graph{}},
		{"one data element", "'data'", network.Graph{}},
		{"one data input", "'data' -> in component(process)", network.Graph{
			Components:  map[string]string{"component": "process"},
			Connections: []network.Connection{DataTarget("data", "in", "component")},
		}},
		{"two connections triple", "'data' -> in one(do) out -> in two(do)", network.Graph{
			Components:  map[string]string{"one": "do", "two": "do"},
			Connections: []network.Connection{DataTarget("data", "in", "one"), SourceTarget("one", "out", "in", "two")},
		}},
		{"two connections separated", "'data' -> in one(do)\none out -> in two(do)", network.Graph{
			Components:  map[string]string{"one": "do", "two": "do"},
			Connections: []network.Connection{DataTarget("data", "in", "one"), SourceTarget("one", "out", "in", "two")},
		}},
		{"split output", "one(do) out -> in two(do)\none out -> in three(do)", network.Graph{
			Components:  map[string]string{"one": "do", "two": "do", "three": "do"},
			Connections: []network.Connection{SourceTarget("one", "out", "in", "two"), SourceTarget("one", "out", "in", "three")},
		}},
		{"join input", "one(do) out -> in three(do)\ntwo(do) out -> in three", network.Graph{
			Components:  map[string]string{"one": "do", "two": "do", "three": "do"},
			Connections: []network.Connection{SourceTarget("one", "out", "in", "three"), SourceTarget("two", "out", "in", "three")},
		}},
		{"index in and out", "one(do) out[1] -> in[1] three(do)\ntwo(do) out[1] -> in[2] three", network.Graph{
			Components: map[string]string{"one": "do", "two": "do", "three": "do"},
			Connections: []network.Connection{
				SourceTargetIndex("one", "out", 1, 1, "in", "three"),
				SourceTargetIndex("two", "out", 1, 2, "in", "three"),
			},
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if g, err := dsl.Parse(strings.NewReader(tc.fbp)); err != nil {
				t.Errorf("Parse failed: %v", err)
			} else if act, exp := g.String(), tc.graph.String(); act != exp {
				t.Errorf("Parse got: %v, want: %v", act, exp)
			}
		})
	}
}

func TestParse_ExampleFile(t *testing.T) {
	if file, err := os.Open("testdata/example.fbp"); err != nil {
		t.Fatal(err)
	} else if graph, err := dsl.Parse(file); err != nil {
		t.Fatal(err)
	} else if len(graph.Components) == 0 {
		t.Error("empty graph")
	}
}
