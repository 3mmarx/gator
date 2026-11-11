package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/3mmarx/gator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("usage: gator login <username>")
	}

	username := cmd.Args[0]

	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		log.Fatalf("user %q not found; run 'gator register %s' to create it", username, username)
	}

	s.cfg.SetUser(username)

	fmt.Printf("Logged in as %q.\n", username)

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("usage: gator register <username>")
	}

	username := cmd.Args[0]

	_, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	})
	if err != nil {
		log.Fatalf("could not register user %q: %v", username, err)
	}

	s.cfg.SetUser(username)

	fmt.Printf("Registered user %q and logged in.\n", username)

	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("could not reset database: %w", err)
	}

	fmt.Println("Database reset. Deleted all users and their data.")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("could not list users: %w", err)
	}

	if len(users) == 0 {
		fmt.Println("No users registered. Run 'gator register <username>'.")
		return nil
	}

	currentUser := s.cfg.CurrentUserName
	for _, user := range users {
		if user.Name == currentUser {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("  %s\n", user.Name)
		}
	}

	return nil
}
