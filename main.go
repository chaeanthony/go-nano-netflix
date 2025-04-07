package main

import (
	"log"
	"net/http"
	"os"

	"github.com/chaeanthony/go-netflix/api"
	"github.com/chaeanthony/go-netflix/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

var DefaultVideoUrl = "https://www.youtube.com/watch?v=ZXsQAXx_ao0"

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

	cfg := api.APIConfig{
		DB:        db,
		JWTSecret: jwtSecret,
		Platform:  platform,
		Port:      port,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/healthz", cfg.HandlerReadiness)

	mux.HandleFunc("POST /auth/signup", cfg.HandlerUsersCreate)
	mux.HandleFunc("POST /auth/login", cfg.HandlerLogin)
	mux.HandleFunc("POST /auth/refresh", cfg.HandlerRefresh)
	mux.HandleFunc("POST /auth/revoke", cfg.HandlerRevoke)

	mux.HandleFunc("GET /api/titles", cfg.HandlerTitlesGet)
	mux.HandleFunc("GET /api/titles/{titleId}", cfg.HandlerTitleGetById)
	mux.HandleFunc("GET /api/shows", cfg.HandlerShowsGet)
	mux.HandleFunc("GET /api/movies", cfg.HandlerMoviesGet)
	mux.Handle("GET /api/watchlist", cfg.AuthTokenMiddleware(http.HandlerFunc(cfg.HandlerWatchlistGet)))
	mux.Handle("POST /api/watchlist", cfg.AuthTokenMiddleware(http.HandlerFunc(cfg.HandlerWatchlistItemCreate)))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on: http://localhost:%s/. Platform: %s\n", port, cfg.Platform)
	log.Fatal(srv.ListenAndServe())
}
