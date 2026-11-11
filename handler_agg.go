package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/3mmarx/gator/internal/database"
	"github.com/google/uuid"
)

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("could not get next feed: %w", err)
	}

	markedFeed, err := s.db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return fmt.Errorf("could not mark feed %q as fetched: %w", feed.Name, err)
	}

	posts, err := fetchFeed(context.Background(), markedFeed.Url)
	if err != nil {
		return fmt.Errorf("could not fetch feed %q: %w", feed.Name, err)
	}

	for _, post := range posts.Channel.Item {
		isValid := true
		pubDate, err := time.Parse(time.RFC1123Z, post.PubDate)
		if err != nil {
			isValid = false
		}

		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       post.Title,
			Url:         post.Link,
			Description: sql.NullString{String: post.Description, Valid: len(post.Description) > 0},
			PublishedAt: sql.NullTime{Time: pubDate, Valid: isValid},
			FeedID:      markedFeed.ID,
		})
		if err != nil {
			// If the URL already exists (a duplicated url)
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Printf("could not store post %q: %v", post.Title, err)
			continue
		}
	}

	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("usage: gator agg <duration> (examples: 30s, 1m, 1h)")
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid duration %q; use a value like 30s, 1m, or 1h", cmd.Args[0])
	}
	if timeBetweenReqs <= 0 {
		return fmt.Errorf("duration must be greater than zero, got %q", cmd.Args[0])
	}

	fmt.Printf("Collecting feeds every %s. Press Ctrl+C to stop.\n", timeBetweenReqs)
	ticker := time.NewTicker(timeBetweenReqs)

	for ; ; <-ticker.C {
		err := scrapeFeeds(s)
		if err != nil {
			return fmt.Errorf("could not scrape feeds: %w", err)
		}
	}
}
