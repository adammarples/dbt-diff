# dbt-diff

A Go CLI tool that helps compare dbt project states between your current branch and `origin/main`, enabling targeted builds and change visualization.

## Installation

```bash
just install
```

This will build and install the binary to `~/.local/bin/dbt-diff`.

## Usage

Must be run from the root of your dbt project (where `dbt_project.yml` exists).

### Build Changed Models

```bash
dbt-diff build-diff
```

This command:
1. Stashes your current changes
2. Fetches and compiles `origin/main` to `target/main/{sha}`
3. Compiles your local changes to `target/local/{diff-hash}`
4. Runs `dbt build --select state:modified` to build only changed models

### Show Changes

```bash
dbt-diff show-diff
```

Displays what models, tests, and other resources have changed compared to `origin/main`, grouped by type and directory.

## How It Works

`dbt-diff` uses dbt's state comparison feature to identify and build only the models that have changed. It:

- Creates unique manifest directories based on git SHAs and diff hashes
- Safely manages git state with automatic stashing and cleanup
- Caches compiled manifests to avoid recompilation
- Provides clear error handling with automatic state restoration

## Development

### Build

```bash
just build
```

Builds the binary to `bin/dbt-diff`.

### Test

```bash
just test
```

Runs all tests.

### Clean

```bash
just clean
```

Removes build artifacts and generated manifest directories.

### Inspect

```bash
just inspect
```

Shows project information, dependencies, and build status.

## Requirements

- Go 1.21+
- `git` in PATH
- `dbt` in PATH
