package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

func establishConnection(addr string) (*websocket.Conn, chan struct{}) {
	c, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				return
			}
			fmt.Print(string(message))
		}
	}()

	return c, done
}

func gracefulClose(c *websocket.Conn, done chan struct{}) {
	err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("close write:", err)
		return
	}
	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
	}
}
