package api

import (
	"context"

	"github.com/walrus811/blog-aggregator/internal/database"
)

func (cfg *ApiConfig) GetPosts(context context.Context, n int) ([]database.Post, error) {
	posts, getPostsErr := cfg.DB.GetPosts(context, int32(n))
	if getPostsErr != nil {
		return nil, getPostsErr
	}

	return posts, nil
}
