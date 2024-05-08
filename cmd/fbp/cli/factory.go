package cli

import (
	"github.com/dgf/go-fbp-x/pkg/process"
	"github.com/dgf/go-fbp-x/pkg/process/core"
	"github.com/dgf/go-fbp-x/pkg/process/filesystem"
	"github.com/dgf/go-fbp-x/pkg/process/html"
	"github.com/dgf/go-fbp-x/pkg/process/http"
	"github.com/dgf/go-fbp-x/pkg/process/text"
)

func withoutMeta(fn func() process.Process) func(map[string]string) process.Process {
	return func(map[string]string) process.Process {
		return fn()
	}
}

func NewFactory(out chan<- string) process.Factory {
	return process.NewFactory(map[string]func(map[string]string) process.Process{
		"core/Clone":     withoutMeta(core.Clone),
		"core/Count":     withoutMeta(core.Count),
		"core/Kick":      withoutMeta(core.Kick),
		"core/Output":    func(map[string]string) process.Process { return core.Output(out) },
		"core/Tick":      withoutMeta(core.Tick),
		"fs/ReadFile":    withoutMeta(filesystem.ReadFile),
		"html/Attribute": withoutMeta(html.Attribute),
		"html/Query":     withoutMeta(html.Query),
		"html/Text":      withoutMeta(html.Text),
		"http/Get":       withoutMeta(http.Get),
		"text/Append":    text.Append,
		"text/Split":     withoutMeta(text.Split),
	})
}
