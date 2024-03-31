package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kong"
	"github.com/dgf/go-fbp-x/cmd/fbp/cli"
	"github.com/dgf/go-fbp-x/web"
)

var (
	Name        = "fbp"
	Version     = "development"
	Description = "Flow Based Programming Executor"
)

type Context struct{}

type ProcsCmd struct{}

func (p *ProcsCmd) Run(_ *Context) error {
	fmt.Println(cli.NewFactory(make(chan string, 1)))
	return nil
}

type RunCmd struct {
	Path  string `arg:"" name:"path" help:"FBP to run." type:"path"`
	Trace bool   `default:"false" help:"Enable trace mode." short:"t"`
}

func (r *RunCmd) Run(_ *Context) error {
	ctx, cancelRun := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	defer close(sigs)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		cancelRun()
	}()

	return cli.Run(ctx, r.Path, r.Trace)
}

type ServeCmd struct {
	Addr string `arg:"" default:":9876" name:"addr" help:"Listen address."`
}

func (s *ServeCmd) Run(_ *Context) error {
	return web.Serve(s.Addr)
}

var cmd struct {
	Procs ProcsCmd `cmd:"" help:"List registred processes."`
	Serve ServeCmd `cmd:"" help:"Serve web interface."`
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
