package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/chaeanthony/go-netflix/internal/auth"
	"github.com/chaeanthony/go-netflix/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if params.Password == "" || params.Email == "" {
		respondError(w, http.StatusBadRequest, "Email and password are required", nil)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	user, err := cfg.db.CreateUser(database.CreateUserParams{
		Email:    params.Email,
		Password: hashedPassword,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	respondJSON(w, http.StatusCreated, user)
}
