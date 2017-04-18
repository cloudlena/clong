package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

func main() {
	hub := newHub()
	go hub.run()

	up := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	http.Handle("/", http.FileServer(http.Dir("public")))
	http.Handle("/screen", ScreenViewHandler())
	http.Handle("/ws/client", ClientConnHandler(hub, up))
	http.Handle("/ws/screen", ScreenConnHandler(hub, up))

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	log.Fatalln(http.ListenAndServe(":"+port, nil))
}
