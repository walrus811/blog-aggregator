package api

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/walrus811/blog-aggregator/internal/database"
)

func (cfg *ApiConfig) GetFeeds(context context.Context) ([]database.Feed, error) {
	feeds, getFeedsErr := cfg.DB.GetFeeds(context)
	if getFeedsErr != nil {
		return nil, getFeedsErr
	}

	return feeds, nil
}

func (cfg *ApiConfig) CreateFeed(context context.Context, name, url string, userId uuid.UUID) (database.Feed, database.FeedFollow, error) {
	feed, createFeedErr := cfg.DB.CreateFeed(context, database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    userId,
	})
	if createFeedErr != nil {
		return database.Feed{}, database.FeedFollow{}, createFeedErr
	}

	feedFollow, createFeedFollowErr := cfg.DB.CreateFeedFollow(context, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    feed.ID,
		UserID:    userId,
	})
	if createFeedFollowErr != nil {
		return database.Feed{}, database.FeedFollow{}, createFeedFollowErr
	}

	return feed, feedFollow, nil
}
