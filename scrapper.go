package main

import (
	"context"
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/2Rahul2/rssagg/internal/database"
	"github.com/google/uuid"
)

func startScrapping(
	db *database.Queries,
	concurrency int,
	timeBetweenRequest time.Duration,
) {
	log.Printf("Scrapping on %v goroutines every %s duration", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)
	// if timeBtwRequest is 1 min ...so every one min channel will execute(for loop)
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedToFetch(
			context.Background(),
			int32(concurrency),
		)

		if err != nil {
			log.Println("Error fetching feeds: ", err)
			continue
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)

			go scrapeFeeds(db, wg, feed)
		}
		wg.Wait()
	}
}

func scrapeFeeds(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()
	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Print("Err: error marking feeds")
		return
	}
	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Print("Err: error reading url")
		return
	}
	for _, item := range rssFeed.Channel.Item {
		pubAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("Err: Could not parse date :%v\n", item.PubDate)
			continue
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Title:     item.Title,
			Description: sql.NullString{
				String: item.Description,
				Valid:  item.Description != "",
			},
			PublishedAt: pubAt,
			Url:         item.Link,
			FeedID:      feed.ID,
		})
		if err != nil {
			log.Print("ERR: Error while saving posts:", err)
		}
	}
	log.Print(feed.Name, len(rssFeed.Channel.Item))
}
