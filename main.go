package main

import (
	"log"
	"net/http"
	"os"

	"github.com/chaeanthony/go-netflix/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

var DefaultVideoUrl = "https://www.youtube.com/watch?v=ZXsQAXx_ao0"

type apiConfig struct {
	db        *database.Client
	jwtSecret string
	platform  string
	port      string
}

func main() {
	godotenv.Load()

	pathToDB := os.Getenv("DB_PATH")
	if pathToDB == "" {
		log.Fatal("DB Url must be set")
	}

	db, err := database.NewClient(pathToDB)
	if err != nil {
		log.Fatalf("Couldn't connect to database: %v", err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM environment variable is not set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	cfg := apiConfig{
		db:        db,
		jwtSecret: jwtSecret,
		platform:  platform,
		port:      port,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	mux.HandleFunc("POST /auth/signup", cfg.handlerUsersCreate)
	mux.HandleFunc("POST /auth/login", cfg.handlerLogin)
	mux.HandleFunc("POST /auth/refresh", cfg.handlerRefresh)
	mux.HandleFunc("POST /auth/revoke", cfg.handlerRevoke)

	mux.HandleFunc("GET /api/titles", cfg.handlerTitlesGet)
	mux.HandleFunc("GET /api/titles/{titleId}", cfg.handlerTitleGetById)
	mux.HandleFunc("GET /api/shows", cfg.handlerShowsGet)
	mux.HandleFunc("GET /api/movies", cfg.handlerMoviesGet)
	mux.Handle("GET /api/watchlist", cfg.AuthTokenMiddleware(http.HandlerFunc(cfg.handlerWatchlistGet)))
	mux.Handle("POST /api/watchlist", cfg.AuthTokenMiddleware(http.HandlerFunc(cfg.handlerWatchlistItemCreate)))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on: http://localhost:%s/. Platform: %s\n", port, cfg.platform)
	log.Fatal(srv.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
