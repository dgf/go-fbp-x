package wasm

import (
	"context"
	"fmt"
	"strings"
	"syscall/js"
	"time"

	"github.com/dgf/go-fbp-x/dsl"
	"github.com/dgf/go-fbp-x/network"
)

func Run(ctx context.Context, id, flow string) error {
	document := js.Global().Get("document")
	console := js.Global().Get("console")

	graph, err := dsl.Parse(strings.NewReader(flow))
	if err != nil {
		return err
	}

	out := make(chan string, 1)
	traces := make(chan network.Trace, 1)
	outElem := document.Call("getElementById", id)

	go func() {
		defer close(out)
		defer close(traces)

		if err := network.NewNetwork(NewFactory(out)).Run(ctx, graph, traces); err != nil {
			div := document.Call("createElement", "div")
			div.Set("innerHTML", fmt.Sprintf("<strong>%v</strong>", err))
			outElem.Call("prepend", div)
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case o := <-out:
				div := document.Call("createElement", "div")
				div.Set("innerHTML", o)
				outElem.Call("prepend", div)
			case t := <-traces:
				console.Call("log", fmt.Sprintf("%v %v %v", time.Now().Format(time.DateTime), id, t.String()))
			}
		}
	}()

	return nil
}
