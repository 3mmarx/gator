package main

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"
)

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("could not list feeds: %w", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found. Add one with 'gator addfeed <name> <url>'.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)

	fmt.Fprintln(w, "NAME\tURL\tOWNER")

	for _, feed := range feeds {
		fmt.Fprintf(w, "%s\t%s\t%s\n",
			feed.Name,
			feed.Url,
			feed.UserName,
		)
	}

	if err := w.Flush(); err != nil {
		return fmt.Errorf("could not write feed table: %w", err)
	}

	return nil
}
