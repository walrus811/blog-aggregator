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

	feedFollow, createFeedFollowErr := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    feed.ID,
		UserID:    user.ID,
	})
	if createFeedFollowErr != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	resObj := CreateFeedResponse{
		Feed:        feed,
		FeedFollows: feedFollow,
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
	Feed        database.Feed       `json:"feed"`
	FeedFollows database.FeedFollow `json:"feed_follows"`
}
