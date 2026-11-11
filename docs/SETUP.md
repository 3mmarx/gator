# Setup

Detailed setup instructions for Gator. For an overview of the project, see the [README](../README.md).

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

### Install with go install

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

## First Run

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

## Development

Build the binary:

```bash
go build -o gator .
./gator browse 10
```

`go run .` is useful while developing. It compiles and runs the program on every invocation.

Regenerate the database layer after changing files in `sql/queries` or `sql/schema`:

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
sqlc generate
```

Run the basic checks:

```bash
go build ./...
go vet ./...
```

## Project Structure

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
├── docs/
│   ├── SETUP.md             # Detailed setup and troubleshooting
│   ├── demo.gif             # Terminal showcase used by the README
│   └── demo.tape            # VHS script for the GIF
├── sqlc.yaml
├── go.mod
└── go.sum
```

## Troubleshooting

| Problem | Fix |
| --- | --- |
| `could not read config file` or `could not parse config file` | Create or fix `~/.gatorconfig.json` as described in [Configuration](#configuration). |
| `connect: connection refused` | Start PostgreSQL, for example `brew services start postgresql@18` or `sudo service postgresql start`. |
| `pq: database "gator" does not exist` | Run `createdb gator`. |
| `pq: relation "users" does not exist` | Apply the migrations with Goose. |
| `zsh: command not found: gator` | Add `$(go env GOPATH)/bin` to your `PATH`. |
| `user "<name>" not found` | Register the user first with `gator register <username>`. |
| `unknown command` from Gator | Check the command name against the [command reference](../README.md#commands). |
| `You are not following any feeds.` or `No posts found.` | Add or follow feeds, then run `gator agg` for a while. |
| `duplicate key value violates unique constraint` when adding a feed | The URL is already registered. Use `gator follow <url>` instead. |
