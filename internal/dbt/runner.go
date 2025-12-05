package dbt

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Runner handles dbt command execution
type Runner struct {
	dbtProjectDir string
}

// New creates a new dbt runner
func New(dbtProjectDir string) *Runner {
	return &Runner{dbtProjectDir: dbtProjectDir}
}

// Model represents a dbt model/test from dbt ls output
type Model struct {
	Name         string `json:"name"`
	ResourceType string `json:"resource_type"`
	PackageName  string `json:"package_name"`
	OriginalPath string `json:"original_file_path"`
	Database     string `json:"database"`
	Schema       string `json:"schema"`
	Alias        string `json:"alias"`
}

// Compile runs dbt compile with the specified target path
func (r *Runner) Compile(targetPath string) error {
	cmd := exec.Command("dbt", "compile", "--target-path", targetPath)
	cmd.Dir = r.dbtProjectDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("dbt compile failed: %w", err)
	}
	return nil
}

// Run runs dbt run with state comparison
func (r *Runner) Run(statePath string) error {
	cmd := exec.Command("dbt", "run", "--select", "state:modified", "--state", statePath)
	cmd.Dir = r.dbtProjectDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("dbt run failed: %w", err)
	}
	return nil
}

// Test runs dbt test with state comparison
func (r *Runner) Test(statePath string) error {
	cmd := exec.Command("dbt", "test", "--select", "state:modified", "--state", statePath)
	cmd.Dir = r.dbtProjectDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("dbt test failed: %w", err)
	}
	return nil
}

// ListModified returns a list of modified resources compared to the state
func (r *Runner) ListModified(statePath string, resourceType string) ([]Model, error) {
	args := []string{"ls", "--select", "state:modified", "--state", statePath, "--output", "json"}
	if resourceType != "" {
		args = append(args, "--resource-type", resourceType)
	}
	cmd := exec.Command("dbt", args...)
	cmd.Dir = r.dbtProjectDir

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("dbt ls failed: %w", err)
	}

	var models []Model
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		var model Model
		if err := json.Unmarshal([]byte(line), &model); err != nil {
			// Skip lines that aren't JSON (like log messages)
			continue
		}
		models = append(models, model)
	}

	return models, nil
}

// CheckDbtInstalled verifies dbt is available
func CheckDbtInstalled() error {
	cmd := exec.Command("dbt", "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("dbt not found in PATH - please install dbt")
	}
	return nil
}
