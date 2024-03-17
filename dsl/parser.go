// parsing FBP Graph DSL, see https://github.com/flowbased/fbp#readme
//
// current limitation:
// - only string IPs
// - component meta data is ignored
// - annotations are removed before processing
// - not supports port index
package dsl

import (
	"fmt"
	"io"
	"regexp"
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

var (
	commentMatch    = regexp.MustCompile("(?m)#.*$")
	connectionSplit = regexp.MustCompile("->")
	lineSplit       = regexp.MustCompile("(?m)$")
	spacesMatch     = regexp.MustCompile("[\t\f\r ]+")
)

func parseDefinition(part string) (in, component, process, out string, err error) {
	paraStart := strings.Index(part, "(")
	paraEnd := strings.Index(part, ")")
	if paraEnd < paraStart {
		err = fmt.Errorf("invalid component definition %q", part)
		return
	}

	portAndName := strings.Split(part[:paraStart], " ")
	if len(portAndName) == 1 {
		component = portAndName[0]
	} else {
		in = portAndName[0]
		component = portAndName[1]
	}

	process = strings.Split(part[paraStart+1:paraEnd], ":")[0] // remove meta
	out = strings.TrimSpace(part[paraEnd+1:])
	return
}

func parse(pos int, part string) (in, component, process, out string, err error) {
	if strings.Contains(part, "(") {
		return parseDefinition(part)
	} else if conn := strings.Split(part, " "); len(conn) < 2 || len(conn) > 3 {
		err = fmt.Errorf("invalid component defintion %q", part)
		return
	} else if len(conn) == 3 {
		return conn[0], conn[1], "", conn[2], nil
	} else {
		if pos == 0 {
			return "", conn[0], "", conn[1], nil
		}
		return conn[0], conn[1], "", "", nil
	}
}

func Parse(src io.Reader) (Graph, error) {
	graph := Graph{
		Components: map[string]string{},
	}

	allBytes, err := io.ReadAll(src)
	if err != nil {
		return graph, err
	}

	withoutComments := commentMatch.ReplaceAll(allBytes, []byte(""))
	joinedSpaces := spacesMatch.ReplaceAll(withoutComments, []byte(" "))
	for _, line := range strings.Split(string(joinedSpaces), "\n") {
		connection := Connection{}
		for p, part := range connectionSplit.Split(line, -1) {
			trimmed := strings.TrimSpace(part)
			if len(trimmed) > 1 { // skip empty lines
				if strings.HasPrefix(trimmed, "'") { // data input
					connection.Data = strings.Trim(trimmed, "'")
				} else if in, component, process, out, err := parse(p, trimmed); err != nil {
					return Graph{}, err
				} else {

					if len(process) > 0 {
						graph.Components[component] = process
					}

					if len(connection.Data) > 0 || len(connection.Source.Component) > 0 {
						connection.Target.Component = component
						connection.Target.Port = strings.ToLower(in) // TODO handle index
						graph.Connections = append(graph.Connections, connection)
						connection = Connection{}
					}

					if len(out) > 0 {
						connection.Source.Component = component
						connection.Source.Port = strings.ToLower(out)
					}
				}
			}
		}
	}

	return graph, nil
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
