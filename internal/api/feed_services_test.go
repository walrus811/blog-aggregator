package api

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestCreateFeeds(t *testing.T) {
	cfg := getTestApiConfig(TEST_ENV_FILE)
	testName := uuid.New().String()
	user, createUserErr := cfg.CreateUser(context.Background(), testName)

	if createUserErr != nil {
		t.Errorf("Error creating user")
	}

	if user.Name != testName {
		t.Errorf("User name does not match")
	}

	type testData struct {
		Name   string
		URL    string
		UserId uuid.UUID
	}

	testFeeds := []testData{
		{"Feed 1", "http://feed1.com", user.ID},
		{"Feed 2", "http://feed2.com", user.ID},
		{"Feed 3", "http://feed3.com", user.ID},
	}

	deleteFeedIds := make([]uuid.UUID, len(testFeeds))

	for i, testFeed := range testFeeds {
		_, _, createFeedErr := cfg.CreateFeed(context.Background(), testFeed.Name, testFeed.URL, testFeed.UserId)
		if createFeedErr != nil {
			t.Error("Error creating feed", createFeedErr)
		}
		deleteFeedIds[i] = testFeed.UserId
	}

	for _, feedId := range deleteFeedIds {
		deleteFeedErr := cfg.DB.DeleteFeedById(context.Background(), feedId)
		if deleteFeedErr != nil {
			t.Error("Error deleting feed", deleteFeedErr)
		}
	}

	deleteUserErr := cfg.DB.DeleteUserById(context.Background(), user.ID)
	if deleteUserErr != nil {
		t.Error("Error deleting user", deleteUserErr)
	}
}

func TestGetFeeds(t *testing.T) {
	cfg := getTestApiConfig(TEST_ENV_FILE)

	originFeeds, georiginFeedsErr := cfg.GetFeeds(context.Background())
	if georiginFeedsErr != nil {
		t.Error("Error getting feeds", georiginFeedsErr)
	}

	testName := uuid.New().String()
	user, createUserErr := cfg.CreateUser(context.Background(), testName)

	if createUserErr != nil {
		t.Errorf("Error creating user")
	}

	if user.Name != testName {
		t.Errorf("User name does not match")
	}

	type testData struct {
		Name   string
		URL    string
		UserId uuid.UUID
	}

	testFeeds := []testData{
		{"Feed 1", "http://feed1.com", user.ID},
		{"Feed 2", "http://feed2.com", user.ID},
		{"Feed 3", "http://feed3.com", user.ID},
	}

	deleteFeedIds := make([]uuid.UUID, len(testFeeds))

	for i, testFeed := range testFeeds {
		_, _, createFeedErr := cfg.CreateFeed(context.Background(), testFeed.Name, testFeed.URL, testFeed.UserId)
		if createFeedErr != nil {
			t.Error("Error creating feed", createFeedErr)
		}
		deleteFeedIds[i] = testFeed.UserId
	}

	feeds, getFeedsErr := cfg.GetFeeds(context.Background())
	if getFeedsErr != nil {
		t.Error("Error getting feeds", getFeedsErr)
	}

	if len(feeds) != len(originFeeds)+len(testFeeds) {
		t.Error("Number of feeds does not match")
	}

	for _, feedId := range deleteFeedIds {
		deleteFeedErr := cfg.DB.DeleteFeedById(context.Background(), feedId)
		if deleteFeedErr != nil {
			t.Error("Error deleting feed", deleteFeedErr)
		}
	}

	deleteUserErr := cfg.DB.DeleteUserById(context.Background(), user.ID)
	if deleteUserErr != nil {
		t.Error("Error deleting user", deleteUserErr)
	}
}
