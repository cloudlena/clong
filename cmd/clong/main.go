package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/mastertinner/clong"
)

func main() {
	port := flag.String("port", "8080", "the port the app should listen on")
	flag.Parse()

	hub := clong.NewHub()
	hub.Run()

	up := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	http.Handle("/", http.FileServer(http.Dir("public")))
	http.Handle("/screen", clong.ScreenViewHandler())
	http.Handle("/scoreboard", clong.ScoreboardViewHandler())
	http.Handle("/ws/controller", clong.ControllerConnHandler(hub, up))
	http.Handle("/ws/screen", clong.ScreenConnHandler(hub, up))

	log.Fatalln(http.ListenAndServe(":"+*port, nil))
}
