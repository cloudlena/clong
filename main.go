// Command clong is a simple game that allows controller-
// and screen devices to connect to each other.
// The goal of the game is to hit targets on the screen
// by flicking balls at them from the controller.
package main

import (
	"database/sql"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cloudlena/adapters/basicauth"
	"github.com/cloudlena/clong/internal/clong"
	"github.com/cloudlena/clong/internal/clong/httpws"
	"github.com/cloudlena/clong/internal/clong/pg"
	_ "github.com/lib/pq"
)

const serverTimeout = 5 * time.Second

// Pages served by their own routes live outside of web/static,
// so the public file server can't serve them without auth.
//
//go:embed web
var webFS embed.FS

func main() {
	port := getenv("PORT", "8080")
	databaseURL := getenv("DATABASE_URL", "postgresql://postgres:clong@?sslmode=disable")
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		log.Fatalln("ADMIN_PASSWORD environment variable must be set")
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("error opening DB connection: %v", err)
	}
	scores, err := pg.NewScoreStore(db)
	if err != nil {
		log.Fatalf("error creating score store: %v", err)
	}
	svc := clong.NewService(scores)

	static, err := fs.Sub(webFS, "web/static")
	if err != nil {
		log.Fatalln(err)
	}

	// Admin endpoints are protected by basic auth
	users := []basicauth.User{{Username: "admin", Password: adminPassword}}

	mux := http.NewServeMux()
	mux.Handle("GET /screen", basicauth.Handler("Clong screen", users)(serveFile(webFS, "web/screen.html")))
	mux.Handle("GET /scoreboard", serveFile(webFS, "web/scoreboard.html"))
	mux.Handle("GET /ws/controller", httpws.HandleControllerConn(svc))
	mux.Handle("GET /ws/screen", httpws.HandleScreenConn(svc))
	mux.Handle("GET /api/scores", httpws.HandleFindScores(scores))
	mux.Handle("DELETE /api/scores", basicauth.Handler("Clong scores", users)(httpws.HandleDeleteScores(scores)))
	mux.Handle("GET /", http.FileServerFS(static))

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  serverTimeout,
		WriteTimeout: serverTimeout,
	}
	log.Fatalln(srv.ListenAndServe())
}

// getenv returns the value of an environment variable or a fallback if it is not set.
func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// serveFile serves a single file from a file system.
func serveFile(fsys fs.FS, name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, fsys, name)
	}
}
