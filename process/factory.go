package process

import (
	"fmt"
	"slices"
	"strings"

	"golang.org/x/exp/maps"
)

type Factory interface {
	Create(name string) (Process, bool)
	Procs() map[string]Process
	String() string
}

type factory struct {
	procs map[string]func() Process
}

func NewFactory(procs map[string]func() Process) Factory {
	return &factory{procs: procs}
}

func (f *factory) Create(name string) (Process, bool) {
	if fn, ok := f.procs[name]; !ok {
		return nil, false
	} else {
		return fn(), true
	}
}

func portsDoc[S fmt.Stringer](ports map[string]S) string {
	sb := strings.Builder{}
	keys := maps.Keys(ports)
	slices.Sort(keys)
	for _, key := range keys {
		sb.WriteString("\t\t")
		sb.WriteString(key)
		sb.WriteString(" ")

		f := ports[key]
		sb.WriteString(f.String())
		sb.WriteString("\n")
	}
	return sb.String()
}

func (f *factory) Procs() map[string]Process {
	procs := map[string]Process{}
	for name, proc := range f.procs {
		procs[name] = proc()
	}
	return procs
}

func (f *factory) String() string {
	sb := strings.Builder{}
	procs := f.Procs()
	names := maps.Keys(procs)
	slices.Sort(names)

	for _, name := range names {
		sb.WriteString("\n")
		sb.WriteString(name)
		sb.WriteString(": ")

		p := procs[name]
		sb.WriteString(p.Description())
		sb.WriteString("\n\tinputs:\n")
		sb.WriteString(portsDoc(p.Inputs()))
		sb.WriteString("\toutputs:\n")
		sb.WriteString(portsDoc(p.Outputs()))
	}

	return sb.String()
}
