package main

import (
	"context"
	"fmt"

	"github.com/3mmarx/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, c command) error {
		if s.cfg.CurrentUserName == "" {
			return fmt.Errorf("no user is logged in; run 'gator register <username>' or 'gator login <username>'")
		}

		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("could not load current user %q; run 'gator login <username>' or 'gator register <username>': %w", s.cfg.CurrentUserName, err)
		}
		return handler(s, c, user)
	}
}
