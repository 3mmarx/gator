# Gator

[![Go Version](https://img.shields.io/github/go-mod/go-version/3mmarx/gator?logo=go&logoColor=white)](https://go.dev/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-14%2B-4169E1?logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

A terminal-based RSS feed aggregator written in Go. Gator lets multiple users add RSS feeds, follow feeds shared by others, collect posts into PostgreSQL, and browse the latest articles from the command line. Built as part of the [Build a Blog Aggregator in Go](https://www.boot.dev/courses/build-blog-aggregator-golang) course on [Boot.dev](https://www.boot.dev/).

![Gator terminal demo showing users, feeds, following, and browse](docs/demo.gif)

## Features

- Multiple local users, each with their own followed feeds
- Add RSS feeds and automatically follow them
- Follow or unfollow feeds added by other users
- Long-running aggregator that fetches the least recently updated feed on a timer
- Posts stored in PostgreSQL and deduplicated by URL
- Browse posts from the feeds you follow, newest first
- Single binary; only PostgreSQL is required at runtime

## Tech Stack

Go · PostgreSQL · sqlc · Goose · lib/pq · google/uuid

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

## How It Works

1. On startup, `main.go` reads `~/.gatorconfig.json`, opens a PostgreSQL connection, and registers every command in a handler map.
2. Commands that need a user are wrapped in `middlewareLoggedIn`, which loads the user named by `current_user_name`.
3. `gator agg` selects the feed with the oldest `last_fetched_at` (or one that has never been fetched), marks it as fetched, downloads its RSS XML, and inserts each item as a post.
4. Post URLs are unique, so running the aggregator repeatedly does not create duplicates.
5. `gator browse` joins `posts` with `feed_follows` for the current user and returns the newest posts first.

## Database Schema

```mermaid
erDiagram
    USERS ||--o{ FEEDS : owns
    USERS ||--o{ FEED_FOLLOWS : follows
    FEEDS ||--o{ FEED_FOLLOWS : "followed by"
    FEEDS ||--o{ POSTS : contains
```

The full schema lives in [`sql/schema`](sql/schema), and the queries used to generate the Go database layer live in [`sql/queries`](sql/queries).

## Running Locally

Gator requires Go 1.27.3+ and PostgreSQL 14+. Setup, configuration, and troubleshooting are in [docs/SETUP.md](docs/SETUP.md).

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
