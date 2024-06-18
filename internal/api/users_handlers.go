package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/walrus811/blog-aggregator/internal/database"
	"github.com/walrus811/blog-aggregator/internal/utils"
)

func (cfg *ApiConfig) HandlerGetUserByApiKey(w http.ResponseWriter, r *http.Request, user database.User) {
	resObj := CraeteUserResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		ApiKey:    user.ApiKey,
	}
	utils.RespondWithJSON(w, http.StatusOK, resObj)
}

func (cfg *ApiConfig) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	reqObj := CreateUserRequest{}
	reqDecodeErr := decoder.Decode(&reqObj)
	if reqDecodeErr != nil {
		log.Fatal(reqDecodeErr)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, createUserErr := cfg.CreateUser(r.Context(), reqObj.Name)
	if createUserErr != nil {
		log.Fatal(reqDecodeErr)
		utils.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	resObj := CraeteUserResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		ApiKey:    user.ApiKey,
	}
	utils.RespondWithJSON(w, http.StatusCreated, resObj)
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
