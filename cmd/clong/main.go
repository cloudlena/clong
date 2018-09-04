// Binary clong is a simple game that allows controller-
// and screen devices to connect to each other.
// The goal of the game is to hit targets on the screen
// by flicking balls at them from the controller.
package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
	"github.com/mastertinner/adapters/basicauth"
	"github.com/mastertinner/adapters/enforcehttps"
	"github.com/mastertinner/clong/internal/app/clong"
	"github.com/mastertinner/clong/internal/app/clong/scores/mysql"
	"github.com/matryer/way"
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
	sess, err := sql.Open("mysql", *dbString)
	if err != nil {
		log.Fatal(errors.Wrap(err, "error opening DB session"))
	}
	db, err := mysql.Make(sess)
	if err != nil {
		log.Fatal(errors.Wrap(err, "error creating DB"))
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

	// Set up router
	r := way.NewRouter()
	r.Handle(http.MethodGet, "/screen", basicauth.Handler("Clong screen", users)(clong.HandleScreenView()))
	r.Handle(http.MethodGet, "/scoreboard", clong.HandleScoreboardView())
	r.Handle(http.MethodGet, "/ws/controller", clong.HandleControllerConn(hub, up))
	r.Handle(http.MethodGet, "/ws/screen", clong.HandleScreenConn(hub, up))
	r.Handle(http.MethodGet, "/api/scores", clong.HandleFindScores(db))
	r.Handle(http.MethodDelete, "/api/scores", basicauth.Handler("Clong scores", users)(clong.HandleDeleteScores(db)))
	r.Handle(http.MethodGet, "/...", http.FileServer(http.Dir("web/static")))

	sr := enforcehttps.Handler(*forceHTTPS)(r)
	log.Fatal(http.ListenAndServe(":"+*port, sr))
}
