//go:build !solution

package main

import (
	"flag"
	"os"
)

var addr = flag.String("addr", "localhost:8080", "address")

func main() {
	flag.Parse()

	stop := setupSignalHandler()
	mes := readInput(os.Stdin)
	c, done := establishConnection(*addr)
	defer c.Close()

	handleMessages(c, done, stop, mes)
}
