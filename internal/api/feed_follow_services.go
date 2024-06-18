package api

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/walrus811/blog-aggregator/internal/database"
)

func (cfg *ApiConfig) GetFeedFollows(context context.Context, userId uuid.UUID) ([]database.FeedFollow, error) {
	feeds, getFeedFollowsErr := cfg.DB.GetFeedFollows(context, userId)
	if getFeedFollowsErr != nil {
		return nil, getFeedFollowsErr
	}

	return feeds, nil
}

func (cfg *ApiConfig) CreateFeedFollow(context context.Context, feedId, userId uuid.UUID) (database.FeedFollow, error) {
	feedFollow, createFeedFollowErr := cfg.DB.CreateFeedFollow(context, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    feedId,
		UserID:    userId,
	})
	if createFeedFollowErr != nil {
		return database.FeedFollow{}, createFeedFollowErr
	}

	return feedFollow, nil
}

func (cfg *ApiConfig) DeleteFeedFollow(context context.Context, feedFollowId, userId uuid.UUID) error {
	deleteFeedFollowErr := cfg.DB.DeleteFeedFollow(context, database.DeleteFeedFollowParams{
		ID:     feedFollowId,
		UserID: userId,
	})
	if deleteFeedFollowErr != nil {
		return deleteFeedFollowErr
	}

	return nil
}
