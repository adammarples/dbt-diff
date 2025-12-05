package cmd

import (
	"github.com/adammarples/dbt-diff/internal/dbt"
	"fmt"
	"os"
)

// Markdown implements the markdown command
func Markdown(dbtOpts dbt.DbtOptions) error {
	// Setup state (compile main and local manifests)
	stateInfo, err := SetupState(dbtOpts)
	if err != nil {
		return err
	}

	workDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	dbtRunner := dbt.NewWithOptions(workDir, dbtOpts)

	// Get modified models only
	models, err := dbtRunner.ListModified(stateInfo.MainManifestPath, "model")
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
