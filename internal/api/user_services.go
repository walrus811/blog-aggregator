package api

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/walrus811/blog-aggregator/internal/database"
)

func (cfg *ApiConfig) CreateUser(context context.Context, userName string) (database.User, error) {
	if len(userName) == 0 {
		return database.User{}, errors.New("name cannot be empty")
	}
	user, createUserErr := cfg.DB.CreateUser(context, database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      userName,
	})

	if createUserErr != nil {
		return database.User{}, createUserErr
	}

	return user, nil
}
