package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/3mmarx/gator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit int32 = 2

	if len(cmd.Args) > 0 {
		n, err := strconv.ParseInt(cmd.Args[0], 10, 32)
		if err != nil || n < 1 {
			return fmt.Errorf("limit must be a positive integer, got %q", cmd.Args[0])
		}
		limit = int32(n)
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	})
	if err != nil {
		return fmt.Errorf("could not fetch posts: %w", err)
	}

	if len(posts) == 0 {
		fmt.Println("No posts found. Run 'gator agg <duration>' to collect posts.")
		return nil
	}

	for i, post := range posts {
		if i > 0 {
			fmt.Println()
		}

		fmt.Printf("Title:       %s\n", post.Title)
		fmt.Printf("URL:         %s\n", post.Url)
		if post.PublishedAt.Valid {
			fmt.Printf("Published:   %s\n", post.PublishedAt.Time.Format("2006-01-02 15:04"))
		}
		if post.Description.Valid {
			fmt.Printf("Description: %s\n", post.Description.String)
		}
	}

	return nil
}
