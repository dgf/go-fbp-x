package dsl

import (
	"log"
	"log/slog"
	"os"

	"github.com/dgf/go-fbp-x/network"
)

func Run(path string, sigs <-chan os.Signal) {
	out := make(chan string, 1)
	exit := make(chan bool, 1)
	traces := make(chan network.Trace, 1)

	if f, err := os.Open(path); err != nil {
		log.Fatalf("Load failed: %v", err)
	} else if g, err := Parse(f); err != nil {
		log.Fatalf("Parse failed: %v", err)
	} else if n, err := network.Create(g, out); err != nil {
		log.Fatalf("Create failed: %v", err)
	} else {
		n.Run(traces)
	}

	go func() {
		for t := range traces {
			slog.Info("trace", "packet", t.Packet, "connection", t.Connection)
		}
	}()

	go func() {
		for {
			select {
			case o := <-out:
				slog.Info("output", "text", o)
			case <-sigs:
				exit <- true
			}
		}
	}()

	<-exit
}