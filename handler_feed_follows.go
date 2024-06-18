package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/walrus811/blog-aggregator/internal/database"
)

func (cfg *apiConfig) handlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, getFeedFollowsErr := cfg.DB.GetFeedFollows(r.Context(), user.ID)
	if getFeedFollowsErr != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	respondWithJSON(w, http.StatusOK, feedFollows)
}

func (cfg *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	decoder := json.NewDecoder(r.Body)
	reqObj := CreateFeedFollowRequest{}
	reqDecodeErr := decoder.Decode(&reqObj)
	if reqDecodeErr != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	feedFollow, createFeedFollowErr := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    reqObj.FeedID,
		UserID:    user.ID,
	})
	if createFeedFollowErr != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	resObj := CreateFeedFollowResponse{
		ID:        feedFollow.ID,
		CreatedAt: feedFollow.CreatedAt,
		UpdatedAt: feedFollow.UpdatedAt,
		FeedId:    feedFollow.FeedID,
		UserId:    feedFollow.UserID,
	}

	respondWithJSON(w, http.StatusCreated, resObj)
}

func (cfg *apiConfig) handlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowID := r.PathValue("feedFollowID")
	if feedFollowID == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid feedFollowID")
		return
	}

	feedFollowUUID, feedFollowUUIDErr := uuid.Parse(feedFollowID)
	if feedFollowUUIDErr != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid feedFollowID")
		return
	}

	deleteFeedFollowErr := cfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowUUID,
		UserID: user.ID,
	})

	if deleteFeedFollowErr != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

type CreateFeedFollowRequest struct {
	FeedID uuid.UUID `json:"feed_id"`
}

type CreateFeedFollowResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	FeedId    uuid.UUID `json:"feed_id"`
	UserId    uuid.UUID `json:"user_id"`
}
