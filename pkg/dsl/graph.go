package dsl

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"golang.org/x/exp/maps"
)

type Link struct {
	Component string
	Port      string
	Index     int
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

func portLabel(l Link) string {
	sb := strings.Builder{}
	sb.WriteString(l.Port)
	if l.Index > 0 {
		sb.WriteString("[")
		sb.WriteString(strconv.Itoa(l.Index))
		sb.WriteString("]")
	}
	return sb.String()
}

func (c Connection) String() string {
	target := c.Target
	if len(c.Data) > 0 {
		return fmt.Sprintf("%s > %s %s", c.Data, portLabel(target), target.Component)
	} else {
		return fmt.Sprintf("%s %s > %s %s", c.Source.Component, portLabel(c.Source), portLabel(target), target.Component)
	}
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
		sb.WriteString(c.String())
		sb.WriteString("\n")
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
