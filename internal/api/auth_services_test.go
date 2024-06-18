package api

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestGetApiKey(t *testing.T) {
	headers := make(map[string][]string)
	headers["Authorization"] = []string{"ApiKey test"}
	apiKey, err := getApiKey(headers)
	if err != nil {
		t.Error("Error getting api key", err)
	}
	if apiKey != "test" {
		t.Error("Api key does not match")
	}
}
func TestGetApiKeyWithErr(t *testing.T) {
	testApiKeys := []string{"", "ApiKey", "ApiKey more than two", "fewjifjewj033j"}
	testHeaders := make([]map[string][]string, len(testApiKeys))
	for i, testApiKey := range testApiKeys {
		header := make(map[string][]string)
		header["Authorization"] = []string{testApiKey}
		testHeaders[i] = header
	}

	for _, testHeader := range testHeaders {
		_, err := getApiKey(testHeader)
		if err == nil {
			t.Error("Error not thrown", err)
		}
	}
}

func TestCheckAuth(t *testing.T) {
	cfg := getTestApiConfig(TEST_ENV_FILE)
	testName := uuid.New().String()
	user, createUserErr := cfg.CreateUser(context.Background(), testName)

	if createUserErr != nil {
		t.Error("Error creating user", createUserErr)
	}

	_, authErr := cfg.CheckAuth(context.Background(), map[string][]string{"Authorization": {"ApiKey " + user.ApiKey}})
	if authErr != nil {
		t.Error("Error checking auth", authErr)
	}

	deleteErr := cfg.DB.DeleteUserById(context.Background(), user.ID)
	if deleteErr != nil {
		t.Error("Error deleting user", deleteErr)
	}
}

func TestCheckAuthWithErr(t *testing.T) {
	cfg := getTestApiConfig(TEST_ENV_FILE)

	_, authErr1 := cfg.CheckAuth(context.Background(), map[string][]string{"Authorization": {""}})
	if authErr1 == nil {
		t.Error("API Key is empty, but error not thrown", authErr1)
	}

	testKey := uuid.New().String()

	_, authErr2 := cfg.CheckAuth(context.Background(), map[string][]string{"Authorization": {"ApiKey " + testKey}})
	if authErr2 == nil {
		t.Error("API Key is invalid, but error not thrown", authErr2)
	}
}
