package httpws

import (
	"log"

	"github.com/gorilla/websocket"
)

var upgrader websocket.Upgrader

// closeConn closes a WebSocket connection and logs any error.
func closeConn(conn *websocket.Conn) {
	if err := conn.Close(); err != nil {
		log.Printf("error closing websocket connection: %v\n", err)
	}
}
