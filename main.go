package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/adammarples/dbt-diff/cmd"
	"github.com/adammarples/dbt-diff/internal/dbt"
)

const version = "0.6.1"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	// Handle help and version commands
	if command == "help" || command == "-h" || command == "--help" {
		printUsage()
		return
	}
	if command == "version" || command == "-v" || command == "--version" {
		fmt.Printf("dbt-diff version %s\n", version)
		return
	}

	// Parse flags for each command
	var dbtOpts dbt.DbtOptions

	switch command {
	case "build":
		buildCmd := flag.NewFlagSet("build", flag.ExitOnError)
		buildCmd.StringVar(&dbtOpts.Target, "target", "", "dbt target environment")
		buildCmd.StringVar(&dbtOpts.Vars, "vars", "", "dbt variables (JSON string)")
		buildCmd.IntVar(&dbtOpts.Threads, "threads", 0, "number of threads for dbt")
		buildCmd.StringVar(&dbtOpts.ProfilesDir, "profiles-dir", "", "dbt profiles directory")
		buildCmd.Parse(os.Args[2:])

		if err := cmd.Build(dbtOpts); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "markdown":
		markdownCmd := flag.NewFlagSet("markdown", flag.ExitOnError)
		markdownCmd.StringVar(&dbtOpts.Target, "target", "", "dbt target environment")
		markdownCmd.StringVar(&dbtOpts.Vars, "vars", "", "dbt variables (JSON string)")
		markdownCmd.IntVar(&dbtOpts.Threads, "threads", 0, "number of threads for dbt")
		markdownCmd.StringVar(&dbtOpts.ProfilesDir, "profiles-dir", "", "dbt profiles directory")
		markdownCmd.Parse(os.Args[2:])

		if err := cmd.Markdown(dbtOpts); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Printf("dbt-diff v%s - Compare and build dbt project changes\n", version)
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  dbt-diff build [flags]       Build models that changed compared to origin/main")
	fmt.Println("  dbt-diff markdown [flags]    Generate SQL snippets for inspecting changed models")
	fmt.Println("  dbt-diff version             Show version information")
	fmt.Println("  dbt-diff help                Show this help message")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --target <name>              dbt target environment (e.g., prod, dev)")
	fmt.Println("  --vars <json>                dbt variables as JSON string")
	fmt.Println("  --threads <n>                number of threads for dbt to use")
	fmt.Println("  --profiles-dir <path>        dbt profiles directory")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  dbt-diff build --target prod")
	fmt.Println("  dbt-diff markdown --target dev --threads 4")
	fmt.Println("  dbt-diff build --vars '{\"key\": \"value\"}'")
	fmt.Println()
	fmt.Println("Must be run from the root of a dbt project (where dbt_project.yml exists)")
}
