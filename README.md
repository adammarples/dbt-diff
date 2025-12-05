# dbt-diff

A Go CLI tool for dbt projects that helps build and visualize only your changed models.

`dbt-diff` efficiently compares your current working code to you main branch by storing manifests in hashed locations inside your target directory for both your main branch HEAD and your working code. Then it uses dbt's parser to decide which models have changed. From here, you can build them, or generate markdown docs describing them for a nice PR description to help your coworkers review them.

The compiled manifests are stored in your target directory under `target/main/<target>/<short_sha>` and `target/local/<target>/<diff_hash>`, where `short_sha` is the short_sha of your main HEAD commit, `diff_sha` is a hash of your git diff from main, and `target` is your dbt profile target.

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
<details>
<summary>Example</summary>

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
</details>

```bash
> dbt-diff markdown
```
A markdown snippet showing the filepaths and the locations of the new models in the database
Perfect for your PR!
<details>
<summary>Example</summary>

### Generate Markdown Snippets for PR descriptions

Sometimes the level of indirection in a dbt project can make it hard to know where your models actually materialized

````bash
> dbt-diff markdown
🔍 Analyzing changes...
🌐 Fetching origin/main...
✅ Using cached main manifest (8508f09)
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

**The `--target` flag changes where manifests are stored:**

```bash
# Without --target (uses "default")
target/main/default/8508f09/
target/local/default/0f035aab/

# With --target prod
target/main/prod/8508f09/
target/local/prod/0f035aab/

# With --target dev
target/main/dev/8508f09/
target/local/dev/0f035aab/
```

This is necessary because different targets compile to different databases/schemas. Running with `--target prod` produces a different manifest than `--target dev`.

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

If target/main/{sha} or target/local/{diff-hash} already exist, no re-compilation is necessary.



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
