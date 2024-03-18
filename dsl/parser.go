// parsing FBP Graph DSL, see https://github.com/flowbased/fbp#readme
//
// current limitation:
// - component meta data is ignored
// - annotations are removed before processing
// - not supports port index
package dsl

import (
	"fmt"
	"io"
	"regexp"
	"strings"
)

var (
	commentMatch    = regexp.MustCompile("(?m)#.*$")
	connectionSplit = regexp.MustCompile("->")
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
