package main

import (
	"context"
    "database/sql"
    "log"
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

	//Loop through rss.Feed.Channel.Items and print each title
	for _, item := range rssFeed.Channel.Items {
		log.Println(item.Title)
	}

}
