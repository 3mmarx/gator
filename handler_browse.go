package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/3mmarx/gator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit int32 = 2

	if len(cmd.Args) > 0 {
		n, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("invalid limit: %w", err)
		}

		limit = int32(n)
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	})
	if err != nil {
		return fmt.Errorf("failed to fetch posts: %w", err)
	}

	for _, post := range posts {
		fmt.Printf("- Title: %s\n", post.Title)
		if post.Description.Valid {
			fmt.Printf("  - Description: %v\n", post.Description.String)
		}

		if post.PublishedAt.Valid {
			fmt.Printf("  - Published At: %v\n", post.PublishedAt.Time.Format(time.UnixDate))
		}
	}

	return nil
}
