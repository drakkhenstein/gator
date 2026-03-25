package main

import (
	"context"
	"time"
	"github.com/google/uuid"
    "database/sql"
    "log"
	"github.com/drakkhenstein/gator/internal/database"
	"github.com/lib/pq"
)

func scrapeFeeds(s *state) {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No feeds to fetch, sleeping...")
			return
		}
		log.Printf("Error getting next feed to fetch: %v", err)
		return
	}

	// Mark the feed as fetched
	err = s.db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Error marking feed %d as fetched: %v", feed.ID, err)
		return
	}
	log.Printf("Successfully scraped feed %d", feed.ID)

	//fetch the RSS feed using URL
	rssFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("Error fetching RSS feed %d: %v", feed.ID, err)
		return
	}
	log.Printf("Fetched RSS feed %d with %d items", feed.ID, len(rssFeed.Channel.Items))

	//Loop through rss.Feed.Channel.Items and save each item to the database
	for _, item := range rssFeed.Channel.Items {
		publishedAt := sql.NullTime{}
		t, parseErr := time.Parse(time.RFC1123Z, item.PubDate)
		if parseErr == nil {
			publishedAt = sql.NullTime{Time: t, Valid: true}
		} else {
			log.Printf("Error parsing publication date for item '%s' in feed %d: %v", item.Title, feed.ID, parseErr)
		}

		_, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title: item.Title,
			Url: item.Link,
			Description: sql.NullString{String: item.Description, Valid: item.Description != ""},
			PublishedAt: publishedAt,
			FeedID: feed.ID,
		})
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
				continue
			}
			log.Printf("Error creating post for feed %d: %v", feed.ID, err)
		}
	}

}
