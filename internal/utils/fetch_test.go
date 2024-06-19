package utils

import (
	"testing"
)

func TestFetchRSSFeeds(t *testing.T) {
	testUrls := []string{
		"https://blog.boot.dev/index.xml",
		"https://wagslane.dev/index.xml",
	}

	for _, url := range testUrls {
		_, fetchErr := FetchRSSFeeds(url)
		if fetchErr != nil {
			t.Errorf("Error fetching RSS feed")
		}
	}
}

func TestFetchRSSFeedsWithErr(t *testing.T) {
	testUrls := []string{
		"abcd3u66frg",
		"https://google.com",
	}

	for _, url := range testUrls {
		_, fetchErr := FetchRSSFeeds(url)
		if fetchErr == nil {
			t.Errorf("tried to fetch invalid RSS feed, but no error was thrown")
		}
	}
}
