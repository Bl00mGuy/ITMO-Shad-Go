package main

import (
	"os"
	"os/signal"
	"syscall"
)

func setupSignalHandler() chan os.Signal {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	return stop
}
