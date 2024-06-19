package workers

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/walrus811/blog-aggregator/internal/database"
	"github.com/walrus811/blog-aggregator/internal/utils"
)

func StartScrape(db *database.Queries, concurrency int, interval time.Duration) {
	log.Printf("Collecting feeds every %s on %v goroutines...", interval, concurrency)
	ticker := time.NewTicker(interval)

	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Println("Error fetching feeds to scrape", err)
			continue
		}
		log.Printf("Found %v feeds to fetch!", len(feeds))

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go ScrapeRSSFeed(db, wg, feed)
		}
		wg.Wait()
	}
}

func ScrapeRSSFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()
	markErr := db.MarkFeedFetched(context.Background(), feed.ID)
	if markErr != nil {
		log.Println("Error marking feed fetched", markErr)
	}

	fetched, fetchErr := utils.FetchRSSFeeds(feed.Url)
	if fetchErr != nil {
		log.Println("Error fetching RSS feed", fetchErr)
		return
	}
	for _, item := range fetched.Channel.Item {
		log.Println("Processing item", item.Title)
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(fetched.Channel.Item))
}
