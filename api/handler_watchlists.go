package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/chaeanthony/go-nano-netflix/utils"
	"github.com/google/uuid"
)

// Gets a user's watchlist. Expects middleware to validate token and provide userID.
func (cfg *APIConfig) HandlerWatchlistGet(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "Couldn't find validate user",
			errors.New("attempted to create watchlist item: couldn't get user id from request context"))
		return
	}

	watchlist, err := cfg.DB.GetWatchlist(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.RespondError(w, http.StatusNotFound, "found no watchlist titles", err)
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, "Couldn't get watchlist", err)
		return
	}

	utils.RespondJSON(w, http.StatusOK, watchlist)
}

// Adds a media title to user's watchlist. Expects middleware to validate token and provide userID.
func (cfg *APIConfig) HandlerWatchlistItemCreate(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "Couldn't find validate user",
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
		utils.RespondError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	watchlistItem, err := cfg.DB.AddWatchlistItem(userID, params.TitleID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Couldn't add watchlist item", err)
		return
	}

	utils.RespondJSON(w, http.StatusOK, watchlistItem)
}
