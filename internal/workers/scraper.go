package workers

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
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
		publishedAt, parseErr := time.Parse(time.RFC1123Z, item.PubDate)
		if parseErr != nil {
			log.Println("Error parsing published at", item.Link, parseErr)
			continue
		}
		_, createPostErr := db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
		})
		if createPostErr != nil {
			switch err := createPostErr.(type) {
			case *pq.Error:
				if err.Code != "23505" {
					log.Println("Error creating post", createPostErr)
				}
			default:
				log.Println("Error creating post", createPostErr)
			}
		}
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(fetched.Channel.Item))
}
