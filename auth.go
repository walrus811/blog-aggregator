package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/walrus811/blog-aggregator/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, getApiKeyErr := getApiKey(r.Header)
		if getApiKeyErr != nil {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		user, getUserErr := cfg.DB.GetUserByApiKey(r.Context(), apiKey)
		if getUserErr != nil {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		handler(w, r, user)
	}
}

func getApiKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("no Authorization header")
	}
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) != 2 || splitAuth[0] != "ApiKey" {
		return "", errors.New("invalid Authorization header")
	}
	return splitAuth[1], nil
}
