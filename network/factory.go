package network

import (
	"slices"
	"strings"

	"github.com/dgf/go-fbp-x/process"
	"github.com/dgf/go-fbp-x/process/core"
	"github.com/dgf/go-fbp-x/process/filesystem"
	"github.com/dgf/go-fbp-x/process/text"
	"golang.org/x/exp/maps"
)

type Factory interface {
	Create(name string) (process.Process, bool)
	String() string
}

type factory struct {
	procs map[string]func() process.Process
}

func NewFactory(out chan<- string) Factory {
	p := map[string]func() process.Process{}

	p["core/Clone"] = func() process.Process { return core.Clone() }
	p["core/Count"] = func() process.Process { return core.Count() }
	p["core/Kick"] = func() process.Process { return core.Kick() }
	p["core/Output"] = func() process.Process { return core.Output(out) }
	p["fs/ReadFile"] = func() process.Process { return filesystem.ReadFile() }
	p["text/Append"] = func() process.Process { return text.Append() }
	p["text/Split"] = func() process.Process { return text.Split() }

	return &factory{procs: p}
}

func (f *factory) Create(name string) (process.Process, bool) {
	if fn, ok := f.procs[name]; !ok {
		return nil, false
	} else {
		return fn(), true
	}
}

func inputStringers(inputs map[string]process.Input) map[string]string {
	s := map[string]string{}
	for n, i := range inputs {
		s[n] = i.String()
	}
	return s
}

func outputStringers(outputs map[string]process.Output) map[string]string {
	s := map[string]string{}
	for n, o := range outputs {
		s[n] = o.String()
	}
	return s
}

func streamDoc(ms map[string]string) string {
	sb := strings.Builder{}
	is := maps.Keys(ms)
	slices.Sort(is)
	for _, name := range is {
		sb.WriteString("\t\t")
		sb.WriteString(name)
		sb.WriteString(" ")

		f := ms[name]
		sb.WriteString(f)
		sb.WriteString("\n")
	}
	return sb.String()
}

func (f *factory) String() string {
	sb := strings.Builder{}
	names := maps.Keys(f.procs)
	slices.Sort(names)

	for _, name := range names {
		sb.WriteString("\n")
		sb.WriteString(name)
		sb.WriteString(": ")

		p := f.procs[name]()
		sb.WriteString(p.Description())
		sb.WriteString("\n\tinputs:\n")
		sb.WriteString(streamDoc(inputStringers(p.Inputs())))
		sb.WriteString("\toutputs:\n")
		sb.WriteString(streamDoc(outputStringers(p.Outputs())))
	}

	return sb.String()
}
