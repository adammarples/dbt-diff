package cmd

import (
	"github.com/adammarples/dbt-diff/internal/dbt"
	"fmt"
	"os"
)

// Build implements the build command
func Build() error {
	// Setup state (compile main and local manifests)
	stateInfo, err := SetupState()
	if err != nil {
		return err
	}

	workDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	dbtRunner := dbt.New(workDir)

	// Run modified models
	fmt.Println("🏗️  Running modified models...")
	if err := dbtRunner.Run(stateInfo.MainManifestPath); err != nil {
		return fmt.Errorf("run failed: %w", err)
	}

	fmt.Println("✅ Models run complete!")

	// Test modified models
	fmt.Println("🧪 Testing modified models...")
	if err := dbtRunner.Test(stateInfo.MainManifestPath); err != nil {
		return fmt.Errorf("test failed: %w", err)
	}

	fmt.Println("✅ Tests complete!")
	return nil
}
