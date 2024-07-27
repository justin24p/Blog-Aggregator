package main

import (
	"context"
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/justin24p/rssAggregator/internal/database"
)

// runs in background as server


func startScrapping (
	db *database.Queries,
	// indicate how much goroutines to use for fetching feeds same time
	concurrency int,
	timeBetweenRequest time.Duration,
) {
	log.Printf("Scraping on %v, goroutines every %s duration", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(
			// global context 
			context.Background(), 
			int32(concurrency),

		)
		if err != nil {
			log.Println("error fetching feeds:",err)
			continue
		} 
		// fetch each feed indvidually at the same time 
		wg := &sync.WaitGroup{} 
		for _, feed := range feeds {
			wg.Add(1)

			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()
	}
}
func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	// waiting for some distinct calls
	// once done scraping feed
	defer wg.Done()
	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error marking feed as fetched:",err)
		return 
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("Error fetching feed:",err)
		return 
	}
	for _, item := range rssFeed.Channel.Item {
		description := sql.NullString{}		
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true 
		}
		pub, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("couldnt parse date %v with err %v", item.PubDate, err)
			continue
		}
		_, err = db.CreatePost(context.Background(),
		database.CreatePostParams{
			ID: uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Title: item.Title,
			Description:  description,
			PublishedAt: pub,
			Url: item.Link,
			FeedID: feed.ID,
		})
		if err != nil {
			log.Println("Failed to create the post:",err)
		}
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}