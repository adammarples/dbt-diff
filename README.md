# dbt-diff

A Go CLI tool for dbt projects that helps build and visualize only your changed models

## Commands

```bash
> dbt-diff build

```

```bash
> dbt-diff markdown

```

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
> dbt-diff build
🔍 Analyzing changes...
📦 Stashing current changes...
🌐 Fetching origin/main...
✅ Using cached main manifest (8508f09)
🔄 Returning to branch feature/add-orders...
📤 Applying stashed changes...
✅ Using cached local manifest (0f035aab)
🏗️  Running modified models...
14:36:43  Running with dbt=1.10.15
14:36:44  Registered adapter: duckdb=1.10.0
14:36:44  Found 2 models, 468 macros
14:36:44  
14:36:44  Concurrency: 1 threads (target='dev')
14:36:44  
14:36:44  1 of 1 START sql table model main.customers .................................... [RUN]
14:36:44  1 of 1 OK created sql table model main.customers ............................... [OK in 0.09s]
14:36:44  
14:36:44  Finished running 1 table model in 0 hours 0 minutes and 0.19 seconds (0.19s).
14:36:44  
14:36:44  Completed successfully
14:36:44  
14:36:44  Done. PASS=1 WARN=0 ERROR=0 SKIP=0 NO-OP=0 TOTAL=1
✅ Models run complete!
🧪 Testing modified models...
14:36:45  Running with dbt=1.10.15
14:36:46  Registered adapter: duckdb=1.10.0
14:36:46  Found 2 models, 468 macros
14:36:46  Nothing to do. Try checking your model configs and model specification args
✅ Tests complete!
```

This command:
1. Stashes your current changes
2. Fetches and compiles `origin/main` to `target/main/{sha}`
3. **Prompts to rebase if your branch is behind origin/main**
4. Compiles your local changes to `target/local/{diff-hash}`
5. Runs `dbt run` on modified models
6. Runs `dbt test` on modified models

If target/main/{sha} or target/local/{diff-hash} already exist, no re-compilation is necessary.

### Generate Markdown Snippets for PR descriptions

Sometimesthe level of indirection in a sdbt project can make it hard to know where your models actually materialized

````bash
> dbt-diff markdown
🔍 Analyzing changes...
📦 Stashing current changes...
🌐 Fetching origin/main...
✅ Using cached main manifest (8508f09)
🔄 Returning to branch feature/add-orders...
📤 Applying stashed changes...
✅ Using cached local manifest (0f035aab)

```sql
-- models/customers.sql
desc table prod.analytics.customers;
select top 10 * from prod.analytics.customers;
```
````

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

## Example Workflow

```bash
# Work on your feature branch
git checkout -b feature/new-models

# Make changes to dbt models
vim models/customers.sql

# Build and test only the changes
dbt-diff build

# Generate SQL snippets for inspection
dbt-diff markdown
```
