package state

import (
	"fmt"
	"os"
	"path/filepath"
)

// Manager handles manifest file locations and validation
type Manager struct {
	projectRoot string
}

// New creates a new state manager
func New(projectRoot string) *Manager {
	return &Manager{projectRoot: projectRoot}
}

// ValidateProjectRoot checks if dbt_project.yml exists
func (m *Manager) ValidateProjectRoot() error {
	dbtProjectFile := filepath.Join(m.projectRoot, "dbt_project.yml")
	if _, err := os.Stat(dbtProjectFile); os.IsNotExist(err) {
		return fmt.Errorf("dbt_project.yml not found in current directory - must run from dbt project root")
	}
	return nil
}

// GetMainManifestPath returns the path to the main branch manifest
func (m *Manager) GetMainManifestPath(shortSha string) string {
	return filepath.Join(m.projectRoot, "target", "main", shortSha)
}

// GetLocalManifestPath returns the path to the local changes manifest
func (m *Manager) GetLocalManifestPath(diffHash string) string {
	return filepath.Join(m.projectRoot, "target", "local", diffHash)
}

// ManifestExists checks if a manifest file exists at the given path
func (m *Manager) ManifestExists(manifestPath string) bool {
	manifestFile := filepath.Join(manifestPath, "manifest.json")
	_, err := os.Stat(manifestFile)
	return err == nil
}

// EnsureTargetDir creates the target directory if it doesn't exist
func (m *Manager) EnsureTargetDir(targetPath string) error {
	return os.MkdirAll(targetPath, 0755)
}

// RemovePartialManifest removes a manifest directory (for cleanup)
func (m *Manager) RemovePartialManifest(manifestPath string) error {
	if _, err := os.Stat(manifestPath); err == nil {
		return os.RemoveAll(manifestPath)
	}
	return nil
}
