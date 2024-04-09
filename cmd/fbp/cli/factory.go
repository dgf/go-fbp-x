package cli

import (
	"github.com/dgf/go-fbp-x/pkg/process"
	"github.com/dgf/go-fbp-x/pkg/process/core"
	"github.com/dgf/go-fbp-x/pkg/process/filesystem"
	"github.com/dgf/go-fbp-x/pkg/process/html"
	"github.com/dgf/go-fbp-x/pkg/process/http"
	"github.com/dgf/go-fbp-x/pkg/process/text"
)

func NewFactory(out chan<- string) process.Factory {
	return process.NewFactory(map[string]func() process.Process{
		"core/Clone":  func() process.Process { return core.Clone() },
		"core/Count":  func() process.Process { return core.Count() },
		"core/Kick":   func() process.Process { return core.Kick() },
		"core/Output": func() process.Process { return core.Output(out) },
		"core/Tick":   func() process.Process { return core.Tick() },
		"fs/ReadFile": func() process.Process { return filesystem.ReadFile() },
		"html/Query":  func() process.Process { return html.Query() },
		"http/Get":    func() process.Process { return http.Get() },
		"text/Append": func() process.Process { return text.Append() },
		"text/Split":  func() process.Process { return text.Split() },
	})
}
