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

	http.Handle("/", adapters.Adapt(
		http.FileServer(http.Dir("public")),
		secure.Handler(*forceHTTPS),
	))
	http.Handle("/scores", adapters.Adapt(
		clong.FindScoresHandler(db),
		secure.Handler(*forceHTTPS),
	))
	http.Handle("/screen", adapters.Adapt(
		clong.ScreenViewHandler(),
		secure.Handler(*forceHTTPS),
	))
	http.Handle("/scoreboard", adapters.Adapt(
		clong.ScoreboardViewHandler(),
		secure.Handler(*forceHTTPS),
	))
	http.Handle("/ws/controller", clong.ControllerConnHandler(hub, up))
	http.Handle("/ws/screen", clong.ScreenConnHandler(hub, up))

	log.Fatalln(http.ListenAndServe(":"+*port, nil))
}
