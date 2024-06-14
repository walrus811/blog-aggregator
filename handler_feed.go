package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/walrus811/blog-aggregator/internal/database"
)

func (cfg *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, getFeedsErr := cfg.DB.GetFeeds(r.Context())
	if getFeedsErr != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	resObj := GetFeedsResponse{
		Result: feeds,
	}

	respondWithJSON(w, http.StatusOK, resObj)
}

func (cfg *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	decoder := json.NewDecoder(r.Body)
	reqObj := CreateFeedRequest{}
	reqDecodeErr := decoder.Decode(&reqObj)
	if reqDecodeErr != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	feed, createFeedErr := cfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      reqObj.Name,
		Url:       reqObj.Url,
		UserID:    user.ID,
	})
	if createFeedErr != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	resObj := CreateFeedResponse{
		ID:        feed.ID,
		CreatedAt: feed.CreatedAt,
		UpdatedAt: feed.UpdatedAt,
		Name:      feed.Name,
		Url:       feed.Url,
		UserId:    feed.UserID,
	}

	respondWithJSON(w, http.StatusCreated, resObj)
}

type GetFeedsResponse struct {
	Result []database.Feed `json:"result"`
}

type CreateFeedRequest struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type CreateFeedResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	UserId    uuid.UUID `json:"user_id"`
}
