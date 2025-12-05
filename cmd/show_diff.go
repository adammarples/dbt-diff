package cmd

import (
	"dbt-diff/internal/dbt"
	"dbt-diff/internal/git"
	"dbt-diff/internal/state"
	"fmt"
	"os"
)

// ShowDiff implements the show-diff command
func ShowDiff() error {
	workDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Initialize components
	stateMgr := state.New(workDir)
	gitOps := git.New(workDir)
	dbtRunner := dbt.New(workDir)

	// Validate environment
	if err := stateMgr.ValidateProjectRoot(); err != nil {
		return err
	}

	if err := dbt.CheckDbtInstalled(); err != nil {
		return err
	}

	// Get main SHA for state path
	mainSha, err := gitOps.GetShortSHA()
	if err != nil {
		return err
	}

	mainManifestPath := stateMgr.GetMainManifestPath(mainSha)

	// Check if manifests exist
	if !stateMgr.ManifestExists(mainManifestPath) {
		return fmt.Errorf("main manifest not found at %s - run 'dbt-diff build-diff' first", mainManifestPath)
	}

	fmt.Println("🔍 Analyzing changes...")

	// Get modified models only
	models, err := dbtRunner.ListModified(mainManifestPath, "model")
	if err != nil {
		return err
	}

	if len(models) == 0 {
		fmt.Println("✅ No modified models detected")
		return nil
	}

	// Generate markdown output
	fmt.Println()
	for _, model := range models {
		// Build location from database.schema.alias, omitting empty parts
		var parts []string
		if model.Database != "" {
			parts = append(parts, model.Database)
		}
		if model.Schema != "" {
			parts = append(parts, model.Schema)
		}
		// Always include alias/name
		if model.Alias != "" {
			parts = append(parts, model.Alias)
		} else {
			parts = append(parts, model.Name)
		}

		location := ""
		for i, part := range parts {
			if i > 0 {
				location += "."
			}
			location += part
		}

		fmt.Println("```sql")
		fmt.Printf("-- %s\n", model.OriginalPath)
		fmt.Printf("desc table %s;\n", location)
		fmt.Printf("select top 10 * from %s;\n", location)
		fmt.Println("```")
		fmt.Println()
	}

	return nil
}
