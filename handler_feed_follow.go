package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/3mmarx/gator/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 1 {
		return errors.New("Usage: gator follow <url>")
	}

	feedUrl := cmd.Args[0]

	feed, err := s.db.GetFeedByURL(context.Background(), feedUrl)
	if err != nil {
		return fmt.Errorf("failed to follow feed: %s", err)
	}

	user, err = s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("failed to follow feed: %s", err)
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to follow feed: %s", err)
	}

	fmt.Printf("you followed \"%s\" feed successfully", feedFollow.FeedName)

	return nil
}
