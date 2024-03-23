package dsl

import (
	"log"
	"log/slog"
	"os"

	"github.com/dgf/go-fbp-x/network"
)

func Run(path string, trace bool, exit <-chan bool) {
	out := make(chan string, 1)
	done := make(chan bool, 1)
	traces := make(chan network.Trace, 1)

	if f, err := os.Open(path); err != nil {
		log.Fatalf("Load failed: %v", err)
	} else if g, err := Parse(f); err != nil {
		log.Fatalf("Parse failed: %v", err)
	} else if err := network.NewNetwork(out).Run(g, traces); err != nil {
		log.Fatalf("Run failed: %v", err)
	}

	go func() {
		for t := range traces {
			if trace {
				slog.Info("trace", "packet", t.Packet, "connection", t.Connection)
			}
		}
	}()

	go func() {
		for {
			select {
			case o := <-out:
				slog.Info("output", "text", o)
			case <-exit:
				done <- true
			}
		}
	}()

	<-done
}
