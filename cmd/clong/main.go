package main

import (
	"flag"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
	. "github.com/mastertinner/clong"
	"github.com/pkg/errors"
)

func main() {
	var (
		port     = flag.String("port", "8080", "the port the app should listen on")
		dbString = flag.String("db-string", "root:pwd@/clong", "the connection string to the DB")
	)
	flag.Parse()

	// Set up DB
	db, err := NewDB(*dbString)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "error creating DB"))
	}
	defer db.Close()

	// Set up WebSocket upgrader
	up := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	// Set up messaging hub
	hub := NewHub(db)
	hub.Run()

	http.Handle("/", http.FileServer(http.Dir("public")))
	http.Handle("/scores", FindScoresHandler(db))
	http.Handle("/screen", ScreenViewHandler())
	http.Handle("/scoreboard", ScoreboardViewHandler())
	http.Handle("/ws/controller", ControllerConnHandler(hub, up))
	http.Handle("/ws/screen", ScreenConnHandler(hub, up))

	log.Fatalln(http.ListenAndServe(":"+*port, nil))
}
