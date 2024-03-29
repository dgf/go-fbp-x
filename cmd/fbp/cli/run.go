package cli

import (
	"context"
	"log"
	"log/slog"
	"os"
	"sync"

	"github.com/dgf/go-fbp-x/dsl"
	"github.com/dgf/go-fbp-x/network"
)

func Run(ctx context.Context, path string, trace bool) {
	wg := sync.WaitGroup{}
	out := make(chan string, 1)
	traces := make(chan network.Trace, 1)

	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Load failed: %v", err)
	}

	graph, err := dsl.Parse(file)
	if err != nil {
		log.Fatalf("Parse failed: %v", err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(traces)
		defer close(out)

		if err := network.NewNetwork(NewFactory(out)).Run(ctx, graph, traces); err != nil {
			log.Fatalf("Run failed: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for t := range traces {
			if trace {
				slog.Info("trace", "packet", t.Packet, "connection", t.Connection)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for o := range out {
			slog.Info("output", "packet", o)
		}
	}()

	wg.Wait()
}
