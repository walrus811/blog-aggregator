package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/walrus811/blog-aggregator/internal/database"
	"github.com/walrus811/blog-aggregator/internal/utils"
)

func (cfg *ApiConfig) HandlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, getFeedsErr := cfg.DB.GetFeeds(r.Context())
	if getFeedsErr != nil {
		log.Fatal(getFeedsErr)
		utils.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	resObj := GetFeedsResponse{
		Result: feeds,
	}

	utils.RespondWithJSON(w, http.StatusOK, resObj)
}

func (cfg *ApiConfig) HandlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	decoder := json.NewDecoder(r.Body)
	reqObj := CreateFeedRequest{}
	reqDecodeErr := decoder.Decode(&reqObj)
	if reqDecodeErr != nil {
		log.Fatal(reqDecodeErr)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	feed, feedFollow, createFeedErr := cfg.CreateFeed(r.Context(), reqObj.Name, reqObj.Url, user.ID)
	if createFeedErr != nil {
		log.Fatal(createFeedErr)
		utils.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	resObj := CreateFeedResponse{
		Feed:        feed,
		FeedFollows: feedFollow,
	}

	utils.RespondWithJSON(w, http.StatusCreated, resObj)
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
