package main

import (
	"net/http"
	"strconv"
	"time"
)

type MediaTitle struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	OriginDate  time.Time `json:"origin_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (cfg *apiConfig) handlerTitlesGet(w http.ResponseWriter, r *http.Request) {
	titles, err := cfg.db.GetTitles()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Couldn't get title", err)
		return
	}

	respondJSON(w, http.StatusOK, titles)
}

func (cfg *apiConfig) handlerTitleGetById(w http.ResponseWriter, r *http.Request) {
	titleIdString := r.PathValue("titleId")
	titleId, err := strconv.Atoi(titleIdString)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to parse title id", err)
		return
	}

	title, err := cfg.db.GetTitleById(titleId)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Couldn't get title", err)
		return
	}

	respondJSON(w, http.StatusOK, title)
}

func (cfg *apiConfig) handlerShowsGet(w http.ResponseWriter, r *http.Request) {
	shows, err := cfg.db.GetShows()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Couldn't get shows", err)
		return
	}

	respondJSON(w, http.StatusOK, shows)
}

func (cfg *apiConfig) handlerMoviesGet(w http.ResponseWriter, r *http.Request) {
	shows, err := cfg.db.GetMovies()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Couldn't get shows", err)
		return
	}

	respondJSON(w, http.StatusOK, shows)
}
