package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

// Gets a user's watchlist. Expects middleware to validate token and provide userID.
func (cfg *apiConfig) handlerWatchlistGet(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)
	if !ok {
		respondError(w, http.StatusUnauthorized, "Couldn't find validate user",
			errors.New("attempted to create watchlist item: couldn't get user id from request context"))
		return
	}

	watchlist, err := cfg.db.GetWatchlist(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusNotFound, "found no watchlist titles", err)
			return
		}
		respondError(w, http.StatusInternalServerError, "Couldn't get watchlist", err)
		return
	}

	respondJSON(w, http.StatusOK, watchlist)
}

// Adds a media title to user's watchlist. Expects middleware to validate token and provide userID.
func (cfg *apiConfig) handlerWatchlistItemCreate(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)
	if !ok {
		respondError(w, http.StatusUnauthorized, "Couldn't find validate user",
			errors.New("attempted to create watchlist item: couldn't get user id from request context"))
		return
	}

	type parameters struct {
		TitleID int `json:"title_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	watchlistItem, err := cfg.db.AddWatchlistItem(userID, params.TitleID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Couldn't add watchlist item", err)
		return
	}

	respondJSON(w, http.StatusOK, watchlistItem)
}
