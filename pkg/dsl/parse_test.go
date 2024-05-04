package dsl_test

import (
	"os"
	"strings"
	"testing"

	"github.com/dgf/go-fbp-x/pkg/dsl"
)

func DataTarget(data, port, component string) dsl.Connection {
	return dsl.Connection{Data: data, Target: dsl.Link{Port: port, Component: component}}
}

func SourceTargetIndex(sComp, sPort string, sIndex, tIndex int, tPort, tComp string) dsl.Connection {
	return dsl.Connection{
		Source: dsl.Link{Port: sPort, Index: sIndex, Component: sComp},
		Target: dsl.Link{Port: tPort, Index: tIndex, Component: tComp},
	}
}

func SourceTarget(sComp, sPort, tPort, tComp string) dsl.Connection {
	return SourceTargetIndex(sComp, sPort, 0, 0, tPort, tComp)
}

func TestParse(t *testing.T) {
	for _, tc := range []struct {
		name  string
		fbp   string
		graph dsl.Graph
	}{
		{"empty", "", dsl.Graph{}},
		{"one comment", "# comment", dsl.Graph{}},
		{"one component", "in call", dsl.Graph{}},
		{"one data element", "'data'", dsl.Graph{}},
		{"one data input", "'data' -> in component( process )", dsl.Graph{
			Components:  map[string]dsl.Process{"component": {Name: "process"}},
			Connections: []dsl.Connection{DataTarget("data", "in", "component")},
		}},
		{"process meta data", "'data' -> in component(process: key = value , num= 123,flag =true)", dsl.Graph{
			Components: map[string]dsl.Process{"component": {
				Name: "process", Meta: map[string]string{"flag": "true", "key": "value", "num": "123"},
			}},
			Connections: []dsl.Connection{DataTarget("data", "in", "component")},
		}},
		{"two connections triple", "'data' -> in one(do) out -> IN two(do)", dsl.Graph{
			Components:  map[string]dsl.Process{"one": {Name: "do"}, "two": {Name: "do"}},
			Connections: []dsl.Connection{DataTarget("data", "in", "one"), SourceTarget("one", "out", "in", "two")},
		}},
		{"two connections separated", "'data' -> IN one(do)\none OUT -> in two(do)", dsl.Graph{
			Components:  map[string]dsl.Process{"one": {Name: "do"}, "two": {Name: "do"}},
			Connections: []dsl.Connection{DataTarget("data", "in", "one"), SourceTarget("one", "out", "in", "two")},
		}},
		{"split output", "one(do) out -> in two(do)\none out -> in three(do)", dsl.Graph{
			Components:  map[string]dsl.Process{"one": {Name: "do"}, "two": {Name: "do"}, "three": {Name: "do"}},
			Connections: []dsl.Connection{SourceTarget("one", "out", "in", "two"), SourceTarget("one", "out", "in", "three")},
		}},
		{"join input", "one(do) out -> in three(do)\ntwo(do) out -> IN three", dsl.Graph{
			Components:  map[string]dsl.Process{"one": {Name: "do"}, "two": {Name: "do"}, "three": {Name: "do"}},
			Connections: []dsl.Connection{SourceTarget("one", "out", "in", "three"), SourceTarget("two", "out", "in", "three")},
		}},
		{"index in and out", "one(do) out[1] -> in[1] three(do)\ntwo(do) out[1] -> in[2] three", dsl.Graph{
			Components: map[string]dsl.Process{"one": {Name: "do"}, "two": {Name: "do"}, "three": {Name: "do"}},
			Connections: []dsl.Connection{
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
