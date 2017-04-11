package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/websocket"
)

// Message is a websocket message sent from a client
type Message struct {
	Color  string  `json:"color"`
	PosX   int     `json:"posX"`
	PosY   int     `json:"posY"`
	SpeedX float64 `json:"speedX"`
	SpeedY float64 `json:"speedY"`
}

// clients are the connected clients
var clients = make(map[*websocket.Conn]bool)

// screens are the connected screens
var screens = make(map[*websocket.Conn]bool)

// broadcast is the channel which handles messages
var broadcast = make(chan Message)

// upgrader upgrades HTTP connections to websocket connections
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// messageHub pushes messages around
func messageHub() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast

		// Send it out to every screen that is currently connected
		for s := range screens {
			err := s.WriteJSON(msg)
			if err != nil {
				s.Close()
				delete(screens, s)
				log.Printf("screen removed (%v connected)", len(screens))
			}
		}
	}
}

func broadcasting(ws *websocket.Conn, screen bool) {
	defer ws.Close()

	for {
		var msg Message

		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error reading message: %s", err.Error())
			if screen {
				delete(screens, ws)
				log.Printf("screen removed (%v connected)", len(screens))
			} else {
				delete(clients, ws)
				log.Printf("client removed (%v connected)", len(clients))
			}
			break
		}

		if !screen {
			broadcast <- msg
		}
	}
}

// handleClientConnections handles websocket connections coming from clients
func handleClientConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		msg := fmt.Sprintf("error upgrading connection: %s", err.Error())
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	clients[ws] = true
	log.Printf("client added. (%v connected)", len(clients))

	go broadcasting(ws, false)
}

// handleScreenConnections handles websocket connections coming from screens
func handleScreenConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		msg := fmt.Sprintf("error upgrading connection: %s", err.Error())
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	screens[ws] = true
	log.Printf("screen added. (%v connected)", len(screens))

	go broadcasting(ws, true)
}

// handleScreen returns the screen HTML page
func handleScreen(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, path.Join("public", "screen.html"))
}

func main() {
	go messageHub()

	http.Handle("/", http.FileServer(http.Dir("public")))
	http.HandleFunc("/screen", handleScreen)
	http.HandleFunc("/ws/client", handleClientConnections)
	http.HandleFunc("/ws/screen", handleScreenConnections)

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	log.Fatalln(http.ListenAndServe(":"+port, nil))
}
