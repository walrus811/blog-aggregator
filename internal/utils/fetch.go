package utils

import (
	"encoding/xml"
	"net/http"
	"time"
)

func FetchRSSFeeds(url string) (*RSS, error) {
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}
	res, resRrr := httpClient.Get(url)
	if resRrr != nil {
		return nil, resRrr
	}
	defer res.Body.Close()

	decoder := xml.NewDecoder(res.Body)
	resObj := RSS{}
	reqDecodeErr := decoder.Decode(&resObj)
	if reqDecodeErr != nil {
		return nil, reqDecodeErr
	}

	return &resObj, nil
}

type RSS struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Title         string `xml:"title"`
	Link          string `xml:"link"`
	Description   string `xml:"description"`
	Generator     string `xml:"generator"`
	Language      string `xml:"language"`
	LastBuildDate string `xml:"lastBuildDate"`
	Item          []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
}
