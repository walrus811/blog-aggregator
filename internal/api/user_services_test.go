package api

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestCreateUser(t *testing.T) {
	cfg := getTestApiConfig(TEST_ENV_FILE)
	testName := uuid.New().String()
	user, createUserErr := cfg.CreateUser(context.Background(), testName)

	if createUserErr != nil {
		t.Error("Error creating user", createUserErr)
	}

	if user.Name != testName {
		t.Error("User name does not match")
	}

	deleteErr := cfg.DB.DeleteUserById(context.Background(), user.ID)
	if deleteErr != nil {
		t.Error("Error deleting user", deleteErr)
	}
}

func TestCreateUserWithErr(t *testing.T) {
	cfg := getTestApiConfig(TEST_ENV_FILE)
	user, createEmptyUserErr := cfg.CreateUser(context.Background(), "")

	if createEmptyUserErr == nil {
		t.Errorf("tried to create user with empty name, but no error was thrown")
	}

	deleteErr := cfg.DB.DeleteUserById(context.Background(), user.ID)
	if deleteErr != nil {
		t.Error("Error deleting user", deleteErr)
	}
}
