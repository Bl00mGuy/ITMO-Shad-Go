package main

import (
	"github.com/gorilla/websocket"
	"log"
	"os"
)

func handleMessages(c *websocket.Conn, done chan struct{}, stop chan os.Signal, mes chan string) {
	for {
		select {
		case <-done:
			return
		case <-stop:
			log.Println("interrupt")
			gracefulClose(c, done)
			return
		case text := <-mes:
			if err := c.WriteMessage(websocket.TextMessage, []byte(text)); err != nil {
				log.Println("write:", err)
				return
			}
		}
	}
}
