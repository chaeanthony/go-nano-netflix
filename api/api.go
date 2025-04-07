package api

import (
	"net/http"

	"github.com/chaeanthony/go-netflix/internal/database"
)

type APIConfig struct {
	DB        *database.Client
	JWTSecret string
	Platform  string
	Port      string
}

func (cfg *APIConfig) HandlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
