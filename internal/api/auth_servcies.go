package api

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/walrus811/blog-aggregator/internal/database"
)

func (cfg *ApiConfig) CheckAuth(context context.Context, headers http.Header) (database.User, error) {
	apiKey, getApiKeyErr := getApiKey(headers)
	if getApiKeyErr != nil {
		return database.User{}, getApiKeyErr
	}

	user, getUserErr := cfg.DB.GetUserByApiKey(context, apiKey)
	if getUserErr != nil {
		return database.User{}, getUserErr
	}

	return user, nil
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
