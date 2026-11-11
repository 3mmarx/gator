package main

import (
	"context"
	"database/sql"
	"errors"
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
		return err
	}

	markedFeed, err := s.db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return err
	}

	posts, err := fetchFeed(context.Background(), markedFeed.Url)
	if err != nil {
		return err
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
			log.Printf("Couldn't create post: %v", err)
			continue
		}
	}

	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return errors.New("usage: gator agg <time_between_reqs> (ex: 1s, 1m, 1h)")
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return errors.New("error: failed to parse the time, please pass a valid one")
	}

	fmt.Printf("Collecting feeds every %s\n", timeBetweenReqs)
	ticker := time.NewTicker(timeBetweenReqs)

	for ; ; <-ticker.C {
		err := scrapeFeeds(s)
		if err != nil {
			return fmt.Errorf("failed to scrape: %w", err)
		}
	}
}
