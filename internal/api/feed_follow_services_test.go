package api

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestCreateFeedFollow(t *testing.T) {
	cfg := getTestApiConfig(TEST_ENV_FILE)
	publisherName := uuid.New().String()
	publisher, createPublisherUserErr := cfg.CreateUser(context.Background(), publisherName)

	if createPublisherUserErr != nil {
		t.Error("Error creating user(p)", createPublisherUserErr)
	}

	if publisher.Name != publisherName {
		t.Error("User(p) name does not match")
	}

	subscriberName := uuid.New().String()
	subscriber, createSubscriberUserErr := cfg.CreateUser(context.Background(), subscriberName)
	if createSubscriberUserErr != nil {
		t.Error("Error creating user(s)", createSubscriberUserErr)
	}

	if subscriber.Name != subscriberName {
		t.Error("User(s) name does not match")
	}

	type testData struct {
		Name   string
		URL    string
		UserId uuid.UUID
	}

	testFeeds := []testData{
		{"Feed 1", "http://feed1.com", publisher.ID},
		{"Feed 2", "http://feed2.com", publisher.ID},
		{"Feed 3", "http://feed3.com", publisher.ID},
	}

	feedIds := make([]uuid.UUID, len(testFeeds))

	for i, testFeed := range testFeeds {
		feed, _, createFeedErr := cfg.CreateFeed(context.Background(), testFeed.Name, testFeed.URL, testFeed.UserId)
		if createFeedErr != nil {
			t.Error("Error creating feed", createFeedErr)
		}
		feedIds[i] = feed.ID
	}

	for _, feedId := range feedIds {
		_, createFeedFollowErr := cfg.CreateFeedFollow(context.Background(), feedId, subscriber.ID)
		if createFeedFollowErr != nil {
			t.Error("Error creating feed follow", createFeedFollowErr)
		}
	}

	for _, feedId := range feedIds {
		deleteFeedErr := cfg.DB.DeleteFeedById(context.Background(), feedId)
		if deleteFeedErr != nil {
			t.Error("Error deleting feed", deleteFeedErr)
		}
	}

	deleteSubscriberUserErr := cfg.DB.DeleteUserById(context.Background(), subscriber.ID)
	if deleteSubscriberUserErr != nil {
		t.Error("Error deleting user(p)", deleteSubscriberUserErr)
	}

	deletePubliserUserErr := cfg.DB.DeleteUserById(context.Background(), publisher.ID)
	if deletePubliserUserErr != nil {
		t.Error("Error deleting user(s)", deletePubliserUserErr)
	}
}

func TestGetFeedFollow(t *testing.T) {
	cfg := getTestApiConfig(TEST_ENV_FILE)
	userName := uuid.New().String()
	user, createUserErr := cfg.CreateUser(context.Background(), userName)

	if createUserErr != nil {
		t.Error("Error creating user", createUserErr)
	}

	if user.Name != userName {
		t.Error("User name does not match")
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

	feedIds := make([]uuid.UUID, len(testFeeds))

	for i, testFeed := range testFeeds {
		feed, _, createFeedErr := cfg.CreateFeed(context.Background(), testFeed.Name, testFeed.URL, testFeed.UserId)
		if createFeedErr != nil {
			t.Error("Error creating feed", createFeedErr)
		}
		feedIds[i] = feed.ID
	}

	feedFollows, getFeedFollowsErr := cfg.GetFeedFollows(context.Background(), user.ID)
	if getFeedFollowsErr != nil {
		t.Error("Error getting feed follows", getFeedFollowsErr)
	}

	if len(feedFollows) != len(testFeeds) {
		t.Error("Feed follow count does not match")
	}

	for _, feedId := range feedIds {
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

func TestDeleteFeedFollow(t *testing.T) {
	cfg := getTestApiConfig(TEST_ENV_FILE)
	publisherName := uuid.New().String()
	publisher, createPublisherUserErr := cfg.CreateUser(context.Background(), publisherName)

	if createPublisherUserErr != nil {
		t.Error("Error creating user(p)", createPublisherUserErr)
	}

	if publisher.Name != publisherName {
		t.Error("User(p) name does not match")
	}

	subscriberName := uuid.New().String()
	subscriber, createSubscriberUserErr := cfg.CreateUser(context.Background(), subscriberName)
	if createSubscriberUserErr != nil {
		t.Error("Error creating user(s)", createSubscriberUserErr)
	}

	if subscriber.Name != subscriberName {
		t.Error("User(s) name does not match")
	}

	type testData struct {
		Name   string
		URL    string
		UserId uuid.UUID
	}

	testFeeds := []testData{
		{"Feed 1", "http://feed1.com", publisher.ID},
		{"Feed 2", "http://feed2.com", publisher.ID},
		{"Feed 3", "http://feed3.com", publisher.ID},
	}

	feedIds := make([]uuid.UUID, len(testFeeds))

	for i, testFeed := range testFeeds {
		feed, _, createFeedErr := cfg.CreateFeed(context.Background(), testFeed.Name, testFeed.URL, testFeed.UserId)
		if createFeedErr != nil {
			t.Error("Error creating feed", createFeedErr)
		}
		feedIds[i] = feed.ID
	}

	feedFollowIds := make([]uuid.UUID, len(feedIds))

	for i, feedId := range feedIds {
		feedFollow, createFeedFollowErr := cfg.CreateFeedFollow(context.Background(), feedId, subscriber.ID)
		if createFeedFollowErr != nil {
			t.Error("Error creating feed follow", createFeedFollowErr)
		}
		feedFollowIds[i] = feedFollow.ID
	}

	for _, feedFollowId := range feedFollowIds {
		deleteFeedFollowErr := cfg.DeleteFeedFollow(context.Background(), feedFollowId, subscriber.ID)
		if deleteFeedFollowErr != nil {
			t.Error("Error deleting feed follow", deleteFeedFollowErr)
		}
	}

	feedFollows, getFeedFollowsErr := cfg.GetFeedFollows(context.Background(), subscriber.ID)
	if getFeedFollowsErr != nil {
		t.Error("Error getting feed follows", getFeedFollowsErr)
	}

	if len(feedFollows) != 0 {
		t.Error("Feed follow count does not match")
	}

	for _, feedId := range feedIds {
		deleteFeedErr := cfg.DB.DeleteFeedById(context.Background(), feedId)
		if deleteFeedErr != nil {
			t.Error("Error deleting feed", deleteFeedErr)
		}
	}

	deleteSubscriberUserErr := cfg.DB.DeleteUserById(context.Background(), subscriber.ID)
	if deleteSubscriberUserErr != nil {
		t.Error("Error deleting user(p)", deleteSubscriberUserErr)
	}

	deletePubliserUserErr := cfg.DB.DeleteUserById(context.Background(), publisher.ID)
	if deletePubliserUserErr != nil {
		t.Error("Error deleting user(s)", deletePubliserUserErr)
	}
}
