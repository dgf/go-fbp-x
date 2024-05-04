// parsing FBP Graph DSL, see https://github.com/flowbased/fbp#readme
//
// current limitation:
// - component meta data is ignored
// - annotations are removed before processing
package dsl

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

var (
	commentMatch    = regexp.MustCompile("(?m)#.*$")
	connectionSplit = regexp.MustCompile("->")
	spacesMatch     = regexp.MustCompile("[\t\f\r ]+")
	indexPortMatch  = regexp.MustCompile(`(\w+)\[(\d+)\]`)
)

func parseProcess(process string) *Process {
	nameAndMeta := strings.Split(process, ":")
	name := strings.TrimSpace(nameAndMeta[0])
	p := &Process{Name: name}

	if len(nameAndMeta) == 2 {
		parts := strings.Split(nameAndMeta[1], ",")
		if len(parts) > 0 {
			p.Meta = make(map[string]string, len(parts))
			for _, part := range parts {
				keyAndValue := strings.Split(part, "=")
				if len(keyAndValue) == 2 {
					p.Meta[strings.TrimSpace(keyAndValue[0])] = strings.TrimSpace(keyAndValue[1])
				}
			}
		}
	}

	return p
}

func parseDefinition(part string) (in, component, out string, process *Process, err error) {
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

	process = parseProcess(part[paraStart+1 : paraEnd])
	out = strings.TrimSpace(part[paraEnd+1:])
	return
}

func parse(pos int, part string) (in, component, out string, process *Process, err error) {
	if strings.Contains(part, "(") {
		return parseDefinition(part)
	} else if conn := strings.Split(part, " "); len(conn) < 2 || len(conn) > 3 {
		err = fmt.Errorf("invalid component defintion %q", part)
		return
	} else if len(conn) == 3 {
		return conn[0], conn[1], conn[2], nil, nil
	} else {
		if pos == 0 {
			return "", conn[0], conn[1], nil, nil
		}
		return conn[0], conn[1], "", nil, nil
	}
}

func link(component, port string) Link {
	portName := port
	index := 0
	if indexPortMatch.Match([]byte(port)) {
		portAndIndex := indexPortMatch.FindStringSubmatch(port)
		portName = portAndIndex[1]
		index, _ = strconv.Atoi(portAndIndex[2]) // hint: number matched by regexp
	}
	return Link{Component: component, Port: strings.ToLower(portName), Index: index}
}

func Parse(src io.Reader) (Graph, error) {
	graph := Graph{
		Components: map[string]Process{},
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
				} else if in, component, out, process, err := parse(p, trimmed); err != nil {
					return Graph{}, err
				} else {

					if process != nil {
						graph.Components[component] = *process
					}

					if len(connection.Data) > 0 || len(connection.Source.Component) > 0 {
						connection.Target = link(component, in)
						graph.Connections = append(graph.Connections, connection)
						connection = Connection{}
					}

					if len(out) > 0 {
						connection.Source = link(component, out)
					}
				}
			}
		}
	}

	return graph, nil
}
