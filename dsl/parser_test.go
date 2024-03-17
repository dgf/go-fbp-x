package dsl_test

import (
	"os"
	"strings"
	"testing"

	"github.com/dgf/go-fbp-x/dsl"
)

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
		{"one data input", "'data' -> in component(process)", dsl.Graph{
			Components: map[string]string{"component": "process"},
			Connections: []dsl.Connection{
				{Data: "data", Target: dsl.Link{Port: "in", Component: "component"}},
			},
		}},
		{"two connections triple", "'data' -> in one(do) out -> in two(do)", dsl.Graph{
			Components: map[string]string{"one": "do", "two": "do"},
			Connections: []dsl.Connection{
				{
					Data:   "data",
					Target: dsl.Link{Port: "in", Component: "one"},
				},
				{
					Source: dsl.Link{Port: "out", Component: "one"},
					Target: dsl.Link{Port: "in", Component: "two"},
				},
			},
		}},
		{"two connections separated", "'data' -> in one(do)\none out -> in two(do)", dsl.Graph{
			Components: map[string]string{"one": "do", "two": "do"},
			Connections: []dsl.Connection{
				{
					Data:   "data",
					Target: dsl.Link{Port: "in", Component: "one"},
				},
				{
					Source: dsl.Link{Port: "out", Component: "one"},
					Target: dsl.Link{Port: "in", Component: "two"},
				},
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
	if file, err := os.Open("example.fbp"); err != nil {
		t.Fatal(err)
	} else if graph, err := dsl.Parse(file); err != nil {
		t.Fatal(err)
	} else if len(graph.Components) == 0 {
		t.Error("empty graph")
	}
}
