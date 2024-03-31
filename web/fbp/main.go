//go:build js && wasm

package main

import (
	"context"
	"fmt"
	"slices"
	"syscall/js"

	"github.com/dgf/go-fbp-x/web/fbp/wasm"
	"golang.org/x/exp/maps"
)

type instance struct {
	ctx    context.Context
	cancel context.CancelFunc
}

var instances = map[string]instance{}

func errorResult(err string) js.Value {
	return js.ValueOf(map[string]any{"error": err})
}

func portList[S fmt.Stringer](ports map[string]S) []any {
	keys := maps.Keys(ports)
	slices.Sort(keys)
	list := make([]any, len(keys))

	for p, port := range keys {
		list[p] = map[string]any{
			"name": port,
			"type": ports[port].String(),
		}
	}

	return list
}

func procsFBP(_ js.Value, _ []js.Value) any {
	procs := wasm.NewFactory(make(chan string)).Procs()
	names := maps.Keys(procs)
	slices.Sort(names)

	procList := []any{}
	for _, name := range names {
		proc := procs[name]

		procList = append(procList, map[string]any{
			"name":    name,
			"desc":    proc.Description(),
			"inputs":  portList(proc.Inputs()),
			"outputs": portList(proc.Outputs()),
		})
	}

	return js.ValueOf(procList)
}

func runFBP(_ js.Value, args []js.Value) any {
	if len(args) != 2 {
		return errorResult("call runFBP(id, flow)")
	}

	id := args[0].String()
	flow := args[1].String()

	ctx, cancel := context.WithCancel(context.Background())
	instances[id] = instance{ctx, cancel}

	if err := wasm.Run(ctx, id, flow); err != nil {
		return errorResult(err.Error())
	}

	return nil
}

func stopFBP(_ js.Value, args []js.Value) any {
	if len(args) != 1 {
		return errorResult("call stopFBP(id)")
	}

	id := args[0].String()
	if i, ok := instances[id]; !ok {
		return errorResult(fmt.Sprintf("flow instance %v not fount", id))
	} else {
		i.cancel()
	}

	return nil
}

func main() { // calls are synced due to exec by single JS thread
	js.Global().Set("procsFBP", js.FuncOf(procsFBP))
	js.Global().Set("runFBP", js.FuncOf(runFBP))
	js.Global().Set("stopFBP", js.FuncOf(stopFBP))
	select {} // keep it alive
}
