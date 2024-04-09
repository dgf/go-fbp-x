package cli

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	"github.com/dgf/go-fbp-x/pkg/dsl"
	"github.com/dgf/go-fbp-x/pkg/network"
)

func Run(ctx context.Context, path string, trace bool) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open failed: %w", err)
	}

	if err := os.Chdir(filepath.Dir(path)); err != nil {
		return fmt.Errorf("chdir failed: %w", err)
	}

	graph, err := dsl.Parse(file)
	if err != nil {
		return fmt.Errorf("parse failed: %w", err)
	}

	var runErr error
	out := make(chan string, 1)
	traces := make(chan network.Trace, 1)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(traces)
		defer close(out)

		if err := network.NewNetwork(NewFactory(out)).Run(ctx, graph, traces); err != nil {
			runErr = fmt.Errorf("run failed: %w", err)
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
	return runErr
}
