// Command clong is a simple game that allows controller-
// and screen devices to connect to each other.
// The goal of the game is to hit targets on the screen
// by flicking balls at them from the controller.
package main

import (
	"database/sql"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/cloudlena/adapters/basicauth"
	"github.com/cloudlena/adapters/enforcehttps"
	"github.com/cloudlena/clong/internal/clong"
	"github.com/cloudlena/clong/internal/clong/httpws"
	"github.com/cloudlena/clong/internal/clong/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
	"github.com/matryer/way"
)

const kiloByte = 1024

const serverTimeout = 5 * time.Second

// Set up static assets
//
//go:embed web/static
var staticFS embed.FS

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
	db, err := sql.Open("mysql", *dbString)
	if err != nil {
		log.Fatal(fmt.Errorf("error opening DB connection: %w", err))
	}
	defer db.Close()

	// Set up WebSocket upgrader
	up := websocket.Upgrader{
		ReadBufferSize:  kiloByte,
		WriteBufferSize: kiloByte,
	}

	// Set up service
	scores, err := mysql.NewScoreStore(db)
	if err != nil {
		log.Fatal(fmt.Errorf("error creating score store: %w", err))
	}
	svc := clong.NewService(scores)

	// Set up basic auth user
	users := []basicauth.User{{Username: *username, Password: *password}}

	// Set up static files
	static, err := fs.Sub(staticFS, "web/static")
	if err != nil {
		log.Fatal(err)
	}

	// Set up router
	r := way.NewRouter()
	r.Handle(http.MethodGet, "/screen", basicauth.Handler("Clong screen", users)(httpws.HandleScreenView()))
	r.Handle(http.MethodGet, "/scoreboard", httpws.HandleScoreboardView())
	r.Handle(http.MethodGet, "/ws/controller", httpws.HandleControllerConn(svc, up))
	r.Handle(http.MethodGet, "/ws/screen", httpws.HandleScreenConn(svc, up))
	r.Handle(http.MethodGet, "/api/scores", httpws.HandleFindScores(scores))
	r.Handle(
		http.MethodDelete,
		"/api/scores",
		basicauth.Handler("Clong scores", users)(httpws.HandleDeleteScores(scores)),
	)
	r.Handle(http.MethodGet, "/...", http.FileServer(http.FS(static)))

	srv := &http.Server{
		Addr:         ":" + *port,
		Handler:      enforcehttps.Handler(*forceHTTPS)(r),
		ReadTimeout:  serverTimeout,
		WriteTimeout: serverTimeout,
	}
	log.Fatal(srv.ListenAndServe())
}
