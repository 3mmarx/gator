package main

import (
	"context"
	"fmt"
	"time"

	"github.com/3mmarx/gator/internal/database"
	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("usage: gator addfeed <name> <url>")
	}

	name := cmd.Args[0]
	feedURL := cmd.Args[1]

	currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("could not load current user %q: %w", s.cfg.CurrentUserName, err)
	}

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       feedURL,
		UserID:    currentUser.ID,
	})
	if err != nil {
		return fmt.Errorf("could not add feed %q: %w", name, err)
	}

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), UserID: currentUser.ID, FeedID: feed.ID})
	if err != nil {
		return fmt.Errorf("could not follow feed %q: %w", name, err)
	}

	fmt.Printf("Added feed %q and followed it.\n", feed.Name)

	return nil
}
