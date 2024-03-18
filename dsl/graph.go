package dsl

import (
	"fmt"
	"slices"
	"strings"

	"golang.org/x/exp/maps"
)

type Link struct {
	Component string
	Port      string
	// Index     int
}

type Connection struct {
	Data   string
	Source Link
	Target Link
}

type Graph struct {
	Components  map[string]string
	Connections []Connection
}

func (g Graph) String() string {
	sb := strings.Builder{}

	sb.WriteString("\ncomponents:\n")
	keys := maps.Keys(g.Components)
	slices.Sort(keys)
	for _, k := range keys {
		sb.WriteString(fmt.Sprintf("%s > %s\n", k, g.Components[k]))
	}

	sb.WriteString("\nconnections:\n")
	for _, c := range g.Connections {
		if len(c.Data) > 0 {
			sb.WriteString(fmt.Sprintf("%s > %s %s\n", c.Data, c.Target.Port, c.Target.Component))
		} else {
			sb.WriteString(fmt.Sprintf("%s %s > %s %s\n", c.Source.Component, c.Source.Port, c.Target.Port, c.Target.Component))
		}
	}

	return sb.String()
}

func Leaves(connections []Connection) []string {
	leaves := map[string]struct{}{}

	for _, c := range connections { // remember all
		leaves[c.Target.Component] = struct{}{}
	}

	for _, c := range connections { // remove sources
		if len(c.Source.Component) > 0 {
			delete(leaves, c.Source.Component)
		}
	}

	return maps.Keys(leaves)
}
