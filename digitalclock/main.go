//go:build !solution

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	port := flag.String("port", "80", "http server port")
	flag.Parse()

	http.HandleFunc("/", handleClockRequest)

	serverAddress := fmt.Sprintf(":%s", *port)
	log.Printf("Server started on %s\n", serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, nil))
}
