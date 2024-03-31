package wasm

import (
	"github.com/dgf/go-fbp-x/process"
	"github.com/dgf/go-fbp-x/process/core"
	"github.com/dgf/go-fbp-x/process/text"
)

func NewFactory(out chan<- string) process.Factory {
	return process.NewFactory(map[string]func() process.Process{
		"core/Clone":  func() process.Process { return core.Clone() },
		"core/Count":  func() process.Process { return core.Count() },
		"core/Kick":   func() process.Process { return core.Kick() },
		"core/Output": func() process.Process { return core.Output(out) },
		"core/Tick":   func() process.Process { return core.Tick() },
		"text/Append": func() process.Process { return text.Append() },
		"text/Split":  func() process.Process { return text.Split() },
	})
}
