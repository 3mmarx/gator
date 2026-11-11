# Gator

[![Go Version](https://img.shields.io/github/go-mod/go-version/3mmarx/gator?logo=go&logoColor=white)](https://go.dev/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-14%2B-4169E1?logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

A terminal-based RSS feed aggregator written in Go. Gator lets multiple users add RSS feeds, follow feeds shared by others, continuously collect posts into PostgreSQL, and browse the latest articles without leaving the command line.

## Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Database Setup](#database-setup)
- [Configuration](#configuration)
- [Quick Start](#quick-start)
- [Commands](#commands)
- [Database Schema](#database-schema)
- [Development](#development)
- [How It Works](#how-it-works)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)
- [Acknowledgements](#acknowledgements)
- [License](#license)

## Features

- Multiple local users, each with their own followed feeds
- Add RSS feeds and automatically follow them
- Follow or unfollow feeds added by other users
- Long-running aggregator that fetches the least recently updated feed on a timer
- Posts stored in PostgreSQL and deduplicated by URL
- Browse posts from the feeds you follow, newest first
- Single self-contained binary after installation

## Tech Stack

- **Go** — CLI, HTTP fetching, and RSS/XML parsing
- **PostgreSQL** — persistent storage
- **sqlc** — type-safe database code generated from SQL
- **Goose** — database migrations
- **lib/pq** — PostgreSQL driver
- **google/uuid** — UUID primary keys

## Prerequisites

To build and run Gator you need:

- [Go 1.27.3 or newer](https://go.dev/dl/)
- [PostgreSQL 14 or newer](https://www.postgresql.org/download/) running locally or accessible over the network

### macOS

```bash
brew install go
brew install postgresql@18
brew services start postgresql@18
```

### Debian / Ubuntu / WSL

Install PostgreSQL with `apt`, and install Go from the [official downloads page](https://go.dev/dl/) if your distribution's version is older than 1.27.3.

```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
sudo service postgresql start
```

Verify both tools are available:

```bash
go version
psql --version
```

> Gator compiles to a standalone binary. You need Go to install or build it, but you do **not** need the Go toolchain to run the resulting `gator` binary.

## Installation

Install the CLI with `go install`:

```bash
go install github.com/3mmarx/gator@latest
```

Go places the binary in `$(go env GOPATH)/bin`. Make sure that directory is on your `PATH`:

```bash
export PATH="$(go env GOPATH)/bin:$PATH"
```

Add that line to your `~/.zshrc`, `~/.bashrc`, or equivalent shell profile to make it permanent.

Check the installation:

```bash
gator
```

Without arguments, Gator prints its usage message.

### Build from source

```bash
git clone https://github.com/3mmarx/gator.git
cd gator
go build -o gator .
```

For a fully static binary:

```bash
CGO_ENABLED=0 go build -o gator .
```

Then run the binary directly:

```bash
./gator users
```

## Database Setup

### 1. Create the database

```bash
createdb gator
```

If your PostgreSQL server requires a specific user or password, connect with `psql` and create the database there instead:

```sql
CREATE DATABASE gator;
```

### 2. Install Goose

Gator's schema is managed with [Goose](https://github.com/pressly/goose) migrations located in `sql/schema`:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

### 3. Apply the migrations

The migration files live in the repository, so run this from a clone of the project:

```bash
git clone https://github.com/3mmarx/gator.git
cd gator

goose -dir sql/schema postgres \
  "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable" up
```

Replace the connection string with your own credentials. On macOS with Homebrew, your local PostgreSQL user is often your system username and has no password:

```bash
goose -dir sql/schema postgres \
  "postgres://$(whoami)@localhost:5432/gator?sslmode=disable" up
```

Useful migration commands:

```bash
# Show the status of every migration
goose -dir sql/schema postgres "<connection-string>" status

# Roll back the most recent migration
goose -dir sql/schema postgres "<connection-string>" down
```

## Configuration

Gator reads its configuration from a JSON file in your home directory:

```text
~/.gatorconfig.json
```

Create it before running any command:

```json
{
  "db_url": "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

| Field | Description |
| --- | --- |
| `db_url` | PostgreSQL connection string. The format is `postgres://user:password@host:port/database?sslmode=disable`. |
| `current_user_name` | The user Gator acts as. `gator register` and `gator login` update this field automatically. |

Examples for common local setups:

```json
{
  "db_url": "postgres://ammar@localhost:5432/gator?sslmode=disable",
  "current_user_name": "ammar"
}
```

```json
{
  "db_url": "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
  "current_user_name": "postgres"
}
```

> `~/.gatorconfig.json` stores database credentials. Do not commit it to version control. Gator does not create this file for you.

## Quick Start

Register a user. This creates the account and makes it the current user:

```bash
gator register ammar
```

Add a few RSS feeds. Adding a feed also follows it automatically:

```bash
gator addfeed "Hacker News" "https://news.ycombinator.com/rss"
gator addfeed "Boot.dev Blog" "https://blog.boot.dev/index.xml"
gator addfeed "TechCrunch" "https://techcrunch.com/feed/"
```

Confirm what you are following:

```bash
gator following
```

Start collecting posts. This runs until you stop it with `Ctrl+C`, so leave it in its own terminal:

```bash
gator agg 1m
```

In another terminal, browse the latest posts from the feeds you follow:

```bash
gator browse 10
```

Without a limit, `browse` shows two posts:

```bash
gator browse
```

## Commands

| Command | Description | Requires login |
| --- | --- | --- |
| `gator register <username>` | Create a new user and set them as the current user | No |
| `gator login <username>` | Switch to an existing user | No |
| `gator users` | List all registered users; the current user is marked with `*` | No |
| `gator reset` | Delete all users and, by cascade, their feeds, follows, and posts | No |
| `gator feeds` | List every feed with its name, URL, and owner | No |
| `gator addfeed <name> <url>` | Add a feed and automatically follow it | Yes |
| `gator follow <url>` | Follow a feed that already exists | Yes |
| `gator following` | List the feeds followed by the current user | Yes |
| `gator unfollow <url>` | Unfollow a feed | Yes |
| `gator agg <duration>` | Continuously fetch the least recently fetched feed every `<duration>` | No |
| `gator browse [limit]` | Show recent posts from followed feeds; defaults to 2 | Yes |

Durations use Go syntax, for example `30s`, `1m`, or `1h`. The `agg` command fetches feeds across all users, not just the current user.

## Database Schema

```mermaid
erDiagram
    USERS ||--o{ FEEDS : owns
    USERS ||--o{ FEED_FOLLOWS : follows
    FEEDS ||--o{ FEED_FOLLOWS : "followed by"
    FEEDS ||--o{ POSTS : contains
```

- `users` — registered accounts with unique names
- `feeds` — RSS feeds, each owned by a user, with a nullable `last_fetched_at`
- `feed_follows` — many-to-many follows between users and feeds
- `posts` — fetched articles with a unique URL for deduplication

The full schema lives in [`sql/schema`](sql/schema), and the queries used to generate the Go database layer live in [`sql/queries`](sql/queries).

## Development

Clone the repository and run Gator with `go run`:

```bash
git clone https://github.com/3mmarx/gator.git
cd gator

go run . register ammar
go run . addfeed "Hacker News" "https://news.ycombinator.com/rss"
go run . following
```

`go run .` is for development only. It compiles and runs the program on every invocation. For production, build a binary:

```bash
go build -o gator .
./gator browse 10
```

### Regenerating the database layer

The `internal/database` package is generated by sqlc. After changing files in `sql/queries` or `sql/schema`, install sqlc and regenerate:

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
sqlc generate
```

### Checks

```bash
go build ./...
go vet ./...
```

### Project structure

```text
.
├── main.go                  # Entry point and command registration
├── commands.go              # Command type, registry, and dispatch
├── middleware_loggedin.go   # Resolves the current user for protected commands
├── rss_feed.go              # HTTP fetching and RSS/XML parsing
├── handler_*.go             # One handler per CLI command
├── internal/
│   ├── config/              # ~/.gatorconfig.json handling
│   └── database/            # sqlc-generated database code
├── sql/
│   ├── schema/              # Goose migrations
│   └── queries/             # SQL queries used by sqlc
├── sqlc.yaml
├── go.mod
└── go.sum
```

## How It Works

1. On startup, `main.go` reads `~/.gatorconfig.json`, opens a PostgreSQL connection, and registers every command in a handler map.
2. Commands that need a user are wrapped in `middlewareLoggedIn`, which loads the user named by `current_user_name`.
3. `gator agg` selects the feed with the oldest `last_fetched_at` (or one that has never been fetched), marks it as fetched, downloads its RSS XML, and inserts each item as a post.
4. Post URLs are unique, so running the aggregator repeatedly does not create duplicates.
5. `gator browse` joins `posts` with `feed_follows` for the current user and returns the newest posts first.

## Troubleshooting

| Problem | Fix |
| --- | --- |
| `Failed to read the config file` | Create `~/.gatorconfig.json` as described in [Configuration](#configuration). |
| `connect: connection refused` | Start PostgreSQL, for example `brew services start postgresql@18` or `sudo service postgresql start`. |
| `pq: database "gator" does not exist` | Run `createdb gator`. |
| `pq: relation "users" does not exist` | Apply the migrations with Goose. |
| `zsh: command not found: gator` | Add `$(go env GOPATH)/bin` to your `PATH`. |
| `failed to login with the user: user not found` | Register the user first with `gator register <username>`. |
| `command not found` from Gator | Check the command name against the [command reference](#commands). |
| `no following feeds` or no posts from `browse` | Add or follow feeds, then run `gator agg` for a while. |
| `duplicate key value violates unique constraint` when adding a feed | The URL is already registered. Use `gator follow <url>` instead. |

## Contributing

Issues and pull requests are welcome. To contribute:

1. Fork the repository.
2. Create a feature branch.
3. Make your changes and run `go build ./...` and `go vet ./...`.
4. Open a pull request with a clear description of the change.

## Acknowledgements

Gator was built as part of the [Boot.dev](https://www.boot.dev/) "Build a Blog Aggregator in Go" guided project.

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
