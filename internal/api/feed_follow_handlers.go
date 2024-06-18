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

func (cfg *ApiConfig) HandlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, getFeedFollowsErr := cfg.GetFeedFollows(r.Context(), user.ID)
	if getFeedFollowsErr != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, feedFollows)
}

func (cfg *ApiConfig) HandlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	decoder := json.NewDecoder(r.Body)
	reqObj := CreateFeedFollowRequest{}
	reqDecodeErr := decoder.Decode(&reqObj)
	if reqDecodeErr != nil {
		log.Fatal(reqDecodeErr)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	feedFollow, createFeedFollowErr := cfg.CreateFeedFollow(r.Context(), reqObj.FeedID, user.ID)
	if createFeedFollowErr != nil {
		log.Fatal(createFeedFollowErr)
		utils.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	resObj := CreateFeedFollowResponse{
		ID:        feedFollow.ID,
		CreatedAt: feedFollow.CreatedAt,
		UpdatedAt: feedFollow.UpdatedAt,
		FeedId:    feedFollow.FeedID,
		UserId:    feedFollow.UserID,
	}

	utils.RespondWithJSON(w, http.StatusCreated, resObj)
}

func (cfg *ApiConfig) HandlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowID := r.PathValue("feedFollowID")
	if feedFollowID == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid feedFollowID")
		return
	}

	feedFollowUUID, feedFollowUUIDErr := uuid.Parse(feedFollowID)
	if feedFollowUUIDErr != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid feedFollowID")
		return
	}

	deleteFeedFollowErr := cfg.DeleteFeedFollow(r.Context(), feedFollowUUID, user.ID)

	if deleteFeedFollowErr != nil {
		log.Fatal(deleteFeedFollowErr)
		utils.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	utils.RespondWithJSON(w, http.StatusNoContent, nil)
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
