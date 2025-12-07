# dbt-diff

A Go CLI tool for dbt projects that helps build and visualize only your changed models.

`dbt-diff` efficiently compares your current working code to your main branch by storing a compiled manifest of your main branch HEAD in a hashed location inside your target directory. Your current working code is compiled on-the-fly by dbt commands. Then it uses dbt's state comparison to decide which models have changed. From here, you can build them, or generate markdown docs describing them for a nice PR description to help your coworkers review them.

The main branch manifest is stored in your target directory under `target/main/<target>/<short_sha>`, where `short_sha` is the short SHA of your main HEAD commit and `target` is your dbt profile target. Your local changes are compiled automatically by dbt when running commands.

#### Currrenty not supported: 
* passing --target-path flags to dbt, assumes working in 'target' dir
* remotes other than 'origin' and main branches other than 'main'
* sql dialects other than snowflake in the markdown

## Installation
<details>
<summary>Installation</summary>

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
</details>

## Commands

```bash
> dbt-diff build
```
Run and test only the models that you've changed from your main branch

```bash
22:20:29  1 of 1 START sql table model main.customers .................................... [RUN]
22:20:29  1 of 1 OK created sql table model main.customers ............................... [OK in 0.08s]
```
<details>
<summary>Example</summary>

### Build Changed Models

```bash
> dbt-diff build
🌐 Fetching origin/main...
📦 Stashing current changes...
📝 Compiling origin/main (8508f09)...
23:53:33  Running with dbt=1.10.15
23:53:33  Registered adapter: duckdb=1.10.0
23:53:33  Unable to do partial parsing because saved manifest not found. Starting full parse.
23:53:34  Found 1 model, 468 macros
23:53:34  
23:53:34  Concurrency: 1 threads (target='dev')
23:53:34  
✅ Main manifest compiled
🔄 Returning to previous branch...
📤 Applying stashed changes...
🏗️  Running modified models...
23:53:35  Running with dbt=1.10.15
23:53:35  Registered adapter: duckdb=1.10.0
23:53:35  Unable to do partial parsing because saved manifest not found. Starting full parse.
23:53:36  Found 2 models, 468 macros
23:53:36  
23:53:36  Concurrency: 1 threads (target='dev')
23:53:36  
23:53:36  1 of 1 START sql table model main.orders ....................................... [RUN]
23:53:36  1 of 1 OK created sql table model main.orders .................................. [OK in 0.04s]
23:53:36  
23:53:36  Finished running 1 table model in 0 hours 0 minutes and 0.10 seconds (0.10s).
23:53:36  
23:53:36  Completed successfully
23:53:36  
23:53:36  Done. PASS=1 WARN=0 ERROR=0 SKIP=0 NO-OP=0 TOTAL=1
✅ Models run complete!
🧪 Testing modified models...
23:53:38  Running with dbt=1.10.15
23:53:38  Registered adapter: duckdb=1.10.0
23:53:38  Found 2 models, 468 macros
23:53:38  Nothing to do. Try checking your model configs and model specification args
✅ Tests complete!
```

This command:
1. Stashes your current changes
2. Fetches and compiles `origin/main` to `target/main/{sha}` (if not already cached)
3. **Prompts to rebase if your branch is behind origin/main**
4. Returns to your branch and applies stashed changes
5. Runs `dbt run` on modified models (dbt compiles your current code on-the-fly)
6. Runs `dbt test` on modified models
</details>

```bash
> dbt-diff markdown
```
Make a pretty markdown snippet showing the filepaths and the locations of the branch models in the database. Useful for review, perfect for your PR description.

```sql
-- models/customers.sql
desc table dev.analytics_adam_marples.customers;
select top 10 * from dev.analytics_adam_marples.customers;
```


<details>
<summary>Example</summary>

### Generate Markdown Snippets for PR descriptions

Sometimes the level of indirection in a dbt project can make it hard to know where your models actually materialized

````bash
> dbt-diff markdown
🌐 Fetching origin/main...
✅ Using cached main manifest (8508f09)

```sql
-- models/orders.sql
desc table orders;
select top 10 * from orders;
```

````

## How It Works

`dbt-diff` uses dbt's state comparison feature to identify and build only the models that have changed. It:

- Compiles and caches the main branch manifest based on git SHA
- Safely manages git state with automatic stashing and cleanup
- Lets dbt compile your current changes on-the-fly during command execution
- Provides rebase prompts when your branch is behind main
- Provides clear error handling with automatic state restoration
</details>

## dbt Arguments

`dbt-diff` supports passing common dbt arguments to the underlying dbt commands. These flags are available on both `build` and `markdown` commands. However, because we are hard coding flags like --select, --state, and so on, we only allow the passing of these whitelisted flags.

### Supported Flags

| Flag | Description | Affects Manifest | Example |
|------|-------------|------------------|---------|
| `--target` | dbt target environment | ✅ Yes | `--target prod` |
| `--vars` | Variables passed to dbt | ⚠️ User responsibility | `--vars '{"key":"value"}'` |
| `--threads` | Number of threads | ❌ No | `--threads 4` |
| `--profiles-dir` | Custom profiles directory | ⚠️ User responsibility | `--profiles-dir ~/.dbt` |

### How Flags Affect Manifest Paths

**The `--target` flag changes where the main branch manifest is stored:**

```bash
# Without --target (uses "default")
target/main/default/8508f09/

# With --target prod
target/main/prod/8508f09/

# With --target dev
target/main/dev/8508f09/
```

This is necessary because different targets compile to different databases/schemas. Running with `--target prod` produces a different manifest than `--target dev`. Your local changes are always compiled on-the-fly by dbt, so they don't need a separate cached manifest.

**Other flags (`--vars`, `--threads`, `--profiles-dir`)** are passed through to dbt commands but **do not affect manifest paths**. This means:

- ⚠️ **`--vars`**: If your vars affect compilation, you must use consistent vars for a given target, or manually clear manifests when changing vars
- ✅ **`--threads`**: Safe to change - only affects runtime performance, not manifest content
- ⚠️ **`--profiles-dir`**: If different profiles have different target configurations, use consistent profiles or manually manage manifests

### Examples

```bash
# Build with production target
dbt-diff build --target prod

# Generate markdown with custom vars
dbt-diff markdown --target dev --vars '{"start_date": "2024-01-01"}'
```

## Usage

Must be run from the root of your dbt project (where `dbt_project.yml` exists).

If `target/main/{sha}` already exists, the main branch won't be re-compiled. Your local changes are always compiled fresh by dbt during command execution.



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
