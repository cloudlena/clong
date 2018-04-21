package main

import (
	"flag"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/mastertinner/adapters/basicauth"
	"github.com/mastertinner/adapters/enforcehttps"
	"github.com/mastertinner/clong/internal/app/clong"
	"github.com/pkg/errors"
)

func main() {
	var (
		port       = flag.String("port", "8080", "the port the app should listen on")
		dbString   = flag.String("db-string", "root@/clong", "DB connection string")
		forceHTTPS = flag.Bool("force-https", false, "redirect all requests to HTTPS")
		username   = flag.String("username", "", "username for admin features")
		password   = flag.String("password", "", "password for admin features")
	)
	flag.Parse()

	// Set up DB
	db, err := clong.NewDB(*dbString)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "error creating DB"))
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Fatal(errors.Wrap(err, "error closing DB"))
		}
	}()

	// Set up WebSocket upgrader
	up := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	// Set up messaging hub
	hub := clong.NewHub(db)
	hub.Run()

	// Set up basic auth user
	users := []basicauth.User{{Username: *username, Password: *password}}

	// Set up mux
	r := mux.NewRouter().StrictSlash(true)
	r.Use(enforcehttps.Handler(*forceHTTPS))

	r.
		Path("/screen").
		Handler(basicauth.Handler("Clong screen", users)(
			clong.ScreenViewHandler(),
		))
	r.
		Path("/scoreboard").
		Handler(clong.ScoreboardViewHandler())

	r.
		Path("/ws/controller").
		Handler(clong.ControllerConnHandler(hub, up))
	r.
		Path("/ws/screen").
		Handler(clong.ScreenConnHandler(hub, up))

	r.
		Methods(http.MethodGet).
		Path("/api/scores").
		Handler(clong.FindScoresHandler(db))
	r.
		Methods(http.MethodDelete).
		Path("/api/scores").
		Handler(basicauth.Handler("Clong scores", users)(
			clong.RemoveScoresHandler(db),
		))

	r.
		PathPrefix("/").
		Handler(http.FileServer(http.Dir("web/static")))

	log.Fatalln(http.ListenAndServe(":"+*port, r))
}
