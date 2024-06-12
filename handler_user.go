package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/walrus811/blog-aggregator/internal/database"
)

func (cfg *apiConfig) handlerGetUserByApiKey(w http.ResponseWriter, r *http.Request) {
	apiKey, getApiKeyErr := getApiKey(r.Header)
	if getApiKeyErr != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	user, err := cfg.DB.GetUserByApiKey(r.Context(), apiKey)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Not Found")
		return
	}

	resObj := CraeteUserResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		ApiKey:    user.ApiKey,
	}
	respondWithJSON(w, http.StatusOK, resObj)
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	reqObj := CreateUserRequest{}
	reqDecodeErr := decoder.Decode(&reqObj)
	if reqDecodeErr != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, createUserErr := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      reqObj.Name,
	})
	if createUserErr != nil {
		log.Printf("Responding with 5XX error: %s", createUserErr)
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	resObj := CraeteUserResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		ApiKey:    user.ApiKey,
	}
	respondWithJSON(w, http.StatusCreated, resObj)
}

type GetUserByApiKeyResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	ApiKey    string    `json:"api_key"`
}

type CreateUserRequest struct {
	Name string `json:"name"`
}

type CraeteUserResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	ApiKey    string    `json:"api_key"`
}

func generateApiKey() (string, error) {
	return generateRandomString(32)
}

func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
