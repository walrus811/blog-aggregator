package main

import (
	"errors"
	"net/http"
	"strings"
)

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
