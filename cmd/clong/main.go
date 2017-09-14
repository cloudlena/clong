package main

import (
	"flag"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
	"github.com/mastertinner/adapters"
	"github.com/mastertinner/adapters/secure"
	"github.com/mastertinner/clong"
	"github.com/pkg/errors"
)

func main() {
	var (
		port       = flag.String("port", "8080", "the port the app should listen on")
		dbString   = flag.String("db-string", "root@/clong", "the connection string to the DB")
		forceHTTPS = flag.Bool("force-https", false, "set to redirect any HTTP requests to HTTPS")
	)
	flag.Parse()

	// Set up DB
	db, err := clong.NewDB(*dbString)
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
	hub := clong.NewHub(db)
	hub.Run()

	// Set up mux
	mux := http.NewServeMux()

	mux.Handle("/scores", clong.FindScoresHandler(db))
	mux.Handle("/screen", clong.ScreenViewHandler())
	mux.Handle("/scoreboard", clong.ScoreboardViewHandler())

	mux.Handle("/ws/controller", clong.ControllerConnHandler(hub, up))
	mux.Handle("/ws/screen", clong.ScreenConnHandler(hub, up))

	mux.Handle("/", http.FileServer(http.Dir("public")))

	log.Fatalln(http.ListenAndServe(":"+*port, adapters.Adapt(
		mux,
		secure.Handler(*forceHTTPS),
	)))
}
