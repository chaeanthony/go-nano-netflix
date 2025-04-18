package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/chaeanthony/go-nano-netflix/internal/auth"
	"github.com/chaeanthony/go-nano-netflix/internal/database"
	"github.com/chaeanthony/go-nano-netflix/utils"
)

func (cfg *APIConfig) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		database.User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.DB.GetUserByEmail(params.Email)
	if err != nil {
		utils.RespondError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.Password)
	if err != nil {
		utils.RespondError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	accessToken, err := auth.MakeJWT(
		user.ID,
		cfg.JWTSecret,
		time.Hour,
	)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Couldn't create access JWT", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Couldn't create refresh token", err)
		return
	}

	_, err = cfg.DB.CreateRefreshToken(database.CreateRefreshTokenParams{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().UTC().Add(time.Hour * 24 * 60),
	})
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Couldn't save refresh token", err)
		return
	}

	utils.RespondJSON(w, http.StatusOK, response{
		User:         user,
		Token:        accessToken,
		RefreshToken: refreshToken,
	})
}

func (cfg *APIConfig) HandlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Couldn't find token in header", err)
		return
	}

	user, err := cfg.DB.GetUserByRefreshToken(refreshToken)
	if err != nil {
		utils.RespondError(w, http.StatusUnauthorized, "Couldn't get user for refresh token", err)
		return
	}

	accessToken, err := auth.MakeJWT(
		user.ID,
		cfg.JWTSecret,
		time.Hour,
	)
	if err != nil {
		utils.RespondError(w, http.StatusUnauthorized, "Couldn't validate token", err)
		return
	}

	utils.RespondJSON(w, http.StatusOK, response{
		Token: accessToken,
	})
}

func (cfg *APIConfig) HandlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Couldn't find token", err)
		return
	}

	err = cfg.DB.RevokeRefreshToken(refreshToken)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Couldn't revoke session", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
