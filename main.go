// Command clong is a simple game that allows controller-
// and screen devices to connect to each other.
// The goal of the game is to hit targets on the screen
// by flicking balls at them from the controller.
package main

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cloudlena/adapters/basicauth"
	"github.com/cloudlena/clong/internal/clong"
	"github.com/cloudlena/clong/internal/clong/httpws"
	"github.com/cloudlena/clong/internal/clong/pg"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
)

const (
	kiloByte      = 1024
	serverTimeout = 5 * time.Second
)

//go:embed web/static
var staticFS embed.FS

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}
	databaseURL, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		databaseURL = "postgresql://postgres:clong@?sslmode=disable"
	}
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	// Set up DB
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalln(fmt.Errorf("error opening DB connection: %w", err))
	}
	defer func() {
		err = db.Close()
		if err != nil {
			log.Fatalln(fmt.Errorf("error closing DB connection: %w", err))
		}
	}()

	// Set up WebSocket upgrader
	up := websocket.Upgrader{
		ReadBufferSize:  kiloByte,
		WriteBufferSize: kiloByte,
	}

	// Set up service
	scores, err := pg.NewScoreStore(db)
	if err != nil {
		log.Fatalln(fmt.Errorf("error creating score store: %w", err))
	}
	svc := clong.NewService(scores)

	// Set up basic auth user for admin endpoints
	users := []basicauth.User{{Username: "admin", Password: adminPassword}}

	// Set up static files
	static, err := fs.Sub(staticFS, "web/static")
	if err != nil {
		log.Fatalln(err)
	}

	// Set up router
	mux := http.NewServeMux()
	mux.Handle("GET /screen", basicauth.Handler("Clong screen", users)(httpws.HandleScreenView()))
	mux.Handle("GET /scoreboard", httpws.HandleScoreboardView())
	mux.Handle("GET /ws/controller", httpws.HandleControllerConn(svc, up))
	mux.Handle("GET /ws/screen", httpws.HandleScreenConn(svc, up))
	mux.Handle("GET /api/scores", httpws.HandleFindScores(scores))
	mux.Handle("DELETE /api/scores", basicauth.Handler("Clong scores", users)(httpws.HandleDeleteScores(scores)))
	mux.Handle("GET /", http.FileServer(http.FS(static)))

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  serverTimeout,
		WriteTimeout: serverTimeout,
	}
	log.Fatalln(srv.ListenAndServe())
}
