package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/3mmarx/gator/internal/config"
	"github.com/3mmarx/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	cfg := config.Read()

	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	dbQueries := database.New(db)

	appState := &state{
		db:  dbQueries,
		cfg: &cfg,
	}

	commands := commands{
		handlers: make(map[string]func(*state, command) error),
	}
	commands.register("users", handlerUsers)
	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)

	commands.register("reset", handlerReset)

	commands.register("agg", handlerAgg)

	commands.register("feeds", handlerFeeds)
	commands.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	commands.register("follow", middlewareLoggedIn(handlerFollow))
	commands.register("following", middlewareLoggedIn(handlerFollowing))
	commands.register("unfollow", middlewareLoggedIn(handlerUnfollow))

	commands.register("browse", middlewareLoggedIn(handlerBrowse))

	args := os.Args
	if len(args) < 2 {
		log.Fatal("Usage: gator <command> [args...]")
	}

	command := command{
		Name: args[1],
		Args: args[2:],
	}

	err = commands.run(appState, command)
	if err != nil {
		log.Fatalf("%v", err)
	}
}
