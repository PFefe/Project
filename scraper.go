package main

import (
	"context"
	"database/sql"
	"github.com/MKefem/rssagg/internal/database"
	"github.com/google/uuid"
	"log"
	"strings"
	"sync"
	"time"
)

func startScrapping(
	db *database.Queries,
	concurrency int,
	timeBetweenRequests time.Duration,
) {
	log.Printf(
		"Scrapping on %v goroutines every %s duration",
		concurrency,
		timeBetweenRequests,
	)
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrency),
		)
		if err != nil {
			log.Printf(
				"Failed to get feeds to fetch: %v",
				err,
			)
			continue
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)

			go scrapeFeed(
				db,
				wg,
				feed,
			)
		}
		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	_, err := db.MarkFeedAsFetched(
		context.Background(),
		feed.ID,
	)
	if err != nil {
		log.Printf(
			"Failed to mark feed as fetched: %v",
			err,
		)
		return
	}
	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Printf(
			"Failed to fetch feed: %v",
			err,
		)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}
		pubAt, err := time.Parse(
			time.RFC1123Z,
			item.PubDate,
		)
		if err != nil {
			log.Printf(
				"Failed to parse time: %v witg error %v",
				item.PubDate,
				err,
			)
			continue
		}

		_, err = db.CreatePost(
			context.Background(),
			database.CreatePostParams{
				ID:          uuid.New(),
				CreatedAt:   time.Now().UTC(),
				UpdatedAt:   time.Now().UTC(),
				Title:       item.Title,
				Description: description,
				PublishedAt: pubAt,
				Url:         item.Link,
				FeedID:      feed.ID,
			},
		)
		if err != nil {
			if strings.Contains(
				err.Error(),
				"duplicate key",
			) {
				continue
			}
			log.Printf(
				"Failed to create post: %v",
				err,
			)
		}
	}
	log.Printf(
		"Feed %s collected, %v posts found",
		feed.Name,
		len(rssFeed.Channel.Item),
	)

}
