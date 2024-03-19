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

	"github.com/dgf/go-fbp-x/network"
)

var (
	commentMatch    = regexp.MustCompile("(?m)#.*$")
	connectionSplit = regexp.MustCompile("->")
	spacesMatch     = regexp.MustCompile("[\t\f\r ]+")
	indexPortMatch  = regexp.MustCompile(`(\w+)\[(\d+)\]`)
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

func link(component, port string) network.Link {
	portName := port
	index := 0
	if indexPortMatch.Match([]byte(port)) {
		portAndIndex := indexPortMatch.FindStringSubmatch(port)
		portName = portAndIndex[1]
		index, _ = strconv.Atoi(portAndIndex[2]) // hint: number matched by regexp
	}
	return network.Link{Component: component, Port: strings.ToLower(portName), Index: index}
}

func Parse(src io.Reader) (network.Graph, error) {
	graph := network.Graph{
		Components: map[string]string{},
	}

	allBytes, err := io.ReadAll(src)
	if err != nil {
		return graph, err
	}

	withoutComments := commentMatch.ReplaceAll(allBytes, []byte(""))
	joinedSpaces := spacesMatch.ReplaceAll(withoutComments, []byte(" "))
	for _, line := range strings.Split(string(joinedSpaces), "\n") {
		connection := network.Connection{}
		for p, part := range connectionSplit.Split(line, -1) {
			trimmed := strings.TrimSpace(part)
			if len(trimmed) > 1 { // skip empty lines
				if strings.HasPrefix(trimmed, "'") { // data input
					connection.Data = strings.Trim(trimmed, "'")
				} else if in, component, process, out, err := parse(p, trimmed); err != nil {
					return network.Graph{}, err
				} else {

					if len(process) > 0 {
						graph.Components[component] = process
					}

					if len(connection.Data) > 0 || len(connection.Source.Component) > 0 {
						connection.Target = link(component, in)
						graph.Connections = append(graph.Connections, connection)
						connection = network.Connection{}
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
