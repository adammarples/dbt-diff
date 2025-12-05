# dbt-diff

A Go CLI tool that helps compare dbt project states between your current branch and `origin/main`, enabling targeted builds and change visualization.

## Installation

### Using go install (recommended)

```bash
go install github.com/adammarples/dbt-diff@latest
```

This installs the binary to `$GOPATH/bin` (usually `~/go/bin`). Make sure this is in your PATH.

### From source

```bash
git clone https://github.com/adammarples/dbt-diff
cd dbt-diff
go build -o dbt-diff
# Move to somewhere in your PATH
mv dbt-diff /usr/local/bin/
```

## Usage

Must be run from the root of your dbt project (where `dbt_project.yml` exists).

### Build Changed Models

```bash
dbt-diff build-diff
```

This command:
1. Stashes your current changes
2. Fetches and compiles `origin/main` to `target/main/{sha}`
3. **Prompts to rebase if your branch is behind origin/main**
4. Compiles your local changes to `target/local/{diff-hash}`
5. Runs `dbt run` on modified models
6. Runs `dbt test` on modified models

### Show Changes

```bash
dbt-diff show-diff
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

# See what changed
dbt-diff show-diff

# Build and test only the changes
dbt-diff build-diff

# If prompted, optionally rebase onto latest main
```

## License

MIT
