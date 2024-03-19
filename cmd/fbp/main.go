package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <my-flow-file.fbp", os.Args[0])
	} else {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		Run(os.Args[1], sigs)
	}
}
