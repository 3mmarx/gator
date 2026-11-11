package main

import (
	"context"
	"fmt"

	"github.com/3mmarx/gator/internal/database"
)

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("usage: gator unfollow <url>")
	}

	feedURL := cmd.Args[0]
	feed, err := s.db.GetFeedByURL(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("could not find feed %q: %w", feedURL, err)
	}

	err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("could not unfollow feed %q: %w", feed.Name, err)
	}

	fmt.Printf("Unfollowed %q.\n", feed.Name)

	return nil
}
