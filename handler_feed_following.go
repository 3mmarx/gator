package main

import (
	"context"
	"fmt"

	"github.com/3mmarx/gator/internal/database"
)

func handlerFollowing(s *state, cmd command, user database.User) error {
	currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("could not load current user %q: %w", s.cfg.CurrentUserName, err)
	}

	followedFeeds, err := s.db.GetFeedFollowsForUser(context.Background(), currentUser.ID)
	if err != nil {
		return fmt.Errorf("could not list followed feeds: %w", err)
	}

	if len(followedFeeds) == 0 {
		fmt.Println("You are not following any feeds.")
		return nil
	}

	for _, followedFeed := range followedFeeds {
		fmt.Println(followedFeed.FeedName)
	}

	return nil
}
