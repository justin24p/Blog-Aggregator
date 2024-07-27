package main

import (
	"context"
	"log"
	"sync"
	"time"

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
		log.Println("Erro marking feed as fetched:",err)
		return 
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("Error fetching feed:",err)
		return 
	}
	for _, item := range rssFeed.Channel.Item {
		log.Println("Found post",item.Title)
	}
}