# dbt-diff

A Go CLI tool that helps compare dbt project states between your current branch and `origin/main`, enabling targeted builds and change visualization.

## Installation

### Download pre-built binary (easiest)

Download the latest release for your platform from the [releases page](https://github.com/adammarples/dbt-diff/releases):

**macOS (Apple Silicon):**
```bash
curl -L https://github.com/adammarples/dbt-diff/releases/latest/download/dbt-diff-darwin-arm64.tar.gz | tar xz
sudo mv dbt-diff-darwin-arm64 /usr/local/bin/dbt-diff
```

**macOS (Intel):**
```bash
curl -L https://github.com/adammarples/dbt-diff/releases/latest/download/dbt-diff-darwin-amd64.tar.gz | tar xz
sudo mv dbt-diff-darwin-amd64 /usr/local/bin/dbt-diff
```

**Linux (AMD64):**
```bash
curl -L https://github.com/adammarples/dbt-diff/releases/latest/download/dbt-diff-linux-amd64.tar.gz | tar xz
sudo mv dbt-diff-linux-amd64 /usr/local/bin/dbt-diff
```

**Windows:**
Download `dbt-diff-windows-amd64.zip` from releases and add to PATH.

### Using go install (if you have Go installed)

```bash
go install github.com/adammarples/dbt-diff@latest
```

### From source

```bash
git clone https://github.com/adammarples/dbt-diff
cd dbt-diff
go build -o dbt-diff
sudo mv dbt-diff /usr/local/bin/
```

## Usage

Must be run from the root of your dbt project (where `dbt_project.yml` exists).

### Build Changed Models

```bash
dbt-diff build
```

This command:
1. Stashes your current changes
2. Fetches and compiles `origin/main` to `target/main/{sha}`
3. **Prompts to rebase if your branch is behind origin/main**
4. Compiles your local changes to `target/local/{diff-hash}`
5. Runs `dbt run` on modified models
6. Runs `dbt test` on modified models

### Generate SQL Snippets

```bash
dbt-diff markdown
```

Displays SQL snippets for inspecting changed models:

```sql
-- models/customers.sql
desc table prod.analytics.customers;
select top 10 * from prod.analytics.customers;
```

Ready to copy/paste into your SQL editor!

## How It Works

`dbt-diff` uses dbt's state comparison feature to identify and build only the models that have changed. It:

- Creates unique manifest directories based on git SHAs and diff hashes
- Safely manages git state with automatic stashing and cleanup
- Caches compiled manifests to avoid recompilation
- Provides rebase prompts when your branch is behind main
- Provides clear error handling with automatic state restoration

## Requirements

- Go 1.21+ (for installation)
- `git` in PATH
- `dbt` in PATH (with appropriate adapter for your warehouse)

## Features

- 🚀 **Fast**: Caches manifests - only compiles when code changes
- 🔒 **Safe**: Auto-stashes changes and restores on errors
- 🔄 **Smart**: Prompts to rebase when behind origin/main
- 📝 **Helpful**: Generates SQL snippets for inspecting changes
- ⚡ **Efficient**: Only builds/tests what actually changed

## Example Workflow

```bash
# Work on your feature branch
git checkout -b feature/new-models

# Make changes to dbt models
vim models/customers.sql

# Generate SQL snippets for inspection
dbt-diff markdown

# Build and test only the changes
dbt-diff build

# If prompted, optionally rebase onto latest main
```

## License

MIT
