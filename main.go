package main

import (
	"dbt-diff/cmd"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "build-diff":
		if err := cmd.BuildDiff(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "show-diff":
		if err := cmd.ShowDiff(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("dbt-diff - Compare and build dbt project changes")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  dbt-diff build-diff    Build models that changed compared to origin/main")
	fmt.Println("  dbt-diff show-diff     Show what models/tests have changed")
	fmt.Println("  dbt-diff help          Show this help message")
	fmt.Println()
	fmt.Println("Must be run from the root of a dbt project (where dbt_project.yml exists)")
}
