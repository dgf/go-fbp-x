package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kong"
	"github.com/dgf/go-fbp-x/cmd/fbp/cli"
)

var (
	Name        = "fbp"
	Version     = "development"
	Description = "Flow Based Programming Executor"
)

type Context struct{}

type ProcsCmd struct{}

func (p *ProcsCmd) Run(ctx *Context) error {
	fmt.Println(cli.NewFactory(make(chan string, 1)))
	return nil
}

type RunCmd struct {
	Path  string `arg:"" name:"path" help:"FBP to run." type:"path"`
	Trace bool   `default:"false" help:"Enable trace mode." short:"t"`
}

func (r *RunCmd) Run(ctx *Context) error {
	exit := make(chan bool, 1)
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		exit <- true
	}()

	cli.Run(r.Path, r.Trace, exit)
	return nil
}

var cmd struct {
	Procs ProcsCmd `cmd:"" help:"List registred processes."`
	Run   RunCmd   `cmd:"" help:"Run process."`
}

func main() {
	ctx := kong.Parse(&cmd,
		kong.Name(Name),
		kong.Description(fmt.Sprintf("%s %s", Description, Version)),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: false,
		}),
	)
	ctx.FatalIfErrorf(ctx.Run(&Context{}))
}
