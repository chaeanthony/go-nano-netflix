package main

import (
	"context"
	"net/http"

	"github.com/chaeanthony/go-netflix/internal/auth"
)

type key int

const userIDKey key = 0

func (cfg *apiConfig) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Couldn't find token", err)
			return
		}

		userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Couldn't validate token", err)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
