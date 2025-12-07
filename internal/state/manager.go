package state

import (
	"fmt"
	"os"
	"path/filepath"
)

type Manager struct {
	dbtProjectDir string
}

func New(dbtProjectDir string) *Manager {
	return &Manager{dbtProjectDir: dbtProjectDir}
}

func (m *Manager) ValidateProjectRoot() error {
	dbtProjectFile := filepath.Join(m.dbtProjectDir, "dbt_project.yml")
	if _, err := os.Stat(dbtProjectFile); os.IsNotExist(err) {
		return fmt.Errorf("dbt_project.yml not found in current directory - must run from dbt project root")
	}
	return nil
}

func (m *Manager) GetMainManifestPath(shortSha string, target string) string {
	if target == "" {
		target = "default"
	}
	return filepath.Join(m.dbtProjectDir, "target", "main", target, shortSha)
}

func (m *Manager) ManifestExists(manifestPath string) bool {
	manifestFile := filepath.Join(manifestPath, "manifest.json")
	_, err := os.Stat(manifestFile)
	return err == nil
}

func (m *Manager) EnsureTargetDir(targetPath string) error {
	return os.MkdirAll(targetPath, 0755)
}

func (m *Manager) RemovePartialManifest(manifestPath string) error {
	if _, err := os.Stat(manifestPath); err == nil {
		return os.RemoveAll(manifestPath)
	}
	return nil
}
