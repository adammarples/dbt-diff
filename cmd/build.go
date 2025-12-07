package cmd

import (
	"fmt"
	"os"

	"github.com/adammarples/dbt-diff/internal/dbt"
)

func Build(dbtOpts dbt.DbtOptions) error {

	stateInfo, err := SetupState(dbtOpts)
	if err != nil {
		return err
	}

	workDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	dbtRunner := dbt.NewWithOptions(workDir, dbtOpts)

	fmt.Println("🏗️  Running modified models...")
	if err := dbtRunner.Run(stateInfo.MainManifestPath); err != nil {
		return fmt.Errorf("run failed: %w", err)
	}

	fmt.Println("✅ Models run complete!")

	fmt.Println("🧪 Testing modified models...")
	if err := dbtRunner.Test(stateInfo.MainManifestPath); err != nil {
		return fmt.Errorf("test failed: %w", err)
	}

	fmt.Println("✅ Tests complete!")
	return nil
}
