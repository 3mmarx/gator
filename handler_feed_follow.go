package main

import (
	"context"
	"fmt"
	"time"

	"github.com/3mmarx/gator/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("usage: gator follow <url>")
	}

	feedURL := cmd.Args[0]

	feed, err := s.db.GetFeedByURL(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("could not find feed %q: %w", feedURL, err)
	}

	user, err = s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("could not load current user %q: %w", s.cfg.CurrentUserName, err)
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("could not follow feed %q: %w", feed.Name, err)
	}

	fmt.Printf("Now following %q.\n", feedFollow.FeedName)

	return nil
}
