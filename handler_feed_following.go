package main

import (
	"context"
	"fmt"

	"github.com/3mmarx/gator/internal/database"
)

func handlerFollowing(s *state, cmd command, user database.User) error {
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("failed to fetch following feeds: %s", err)
	}

	followedFeeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("failed to fetch following feeds: %s", err)
	}

	if len(followedFeeds) == 0 {
		fmt.Println("no following feeds")
		return nil
	}

	for _, followedFeed := range followedFeeds {
		fmt.Println(followedFeed.FeedName)
	}

	return nil
}
