package git

import (
	"fmt"
	"os/exec"
	"strings"
)

type Operations struct {
	dbtProjectDir string
}

func New(dbtProjectDir string) *Operations {
	return &Operations{dbtProjectDir: dbtProjectDir}
}

func (g *Operations) GetShortSHA() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--short", "origin/main")
	cmd.Dir = g.dbtProjectDir
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get short SHA: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

func (g *Operations) CreateStash() error {
	cmd := exec.Command("git", "stash", "push")
	cmd.Dir = g.dbtProjectDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create stash: %w", err)
	}
	return nil
}

func (g *Operations) PopStash() error {
	cmd := exec.Command("git", "stash", "pop")
	cmd.Dir = g.dbtProjectDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to pop stash: %w", err)
	}
	return nil
}

func (g *Operations) FetchOrigin() error {
	cmd := exec.Command("git", "fetch", "origin", "main")
	cmd.Dir = g.dbtProjectDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to fetch origin: %w", err)
	}
	return nil
}

func (g *Operations) Checkout(ref string) error {
	cmd := exec.Command("git", "checkout", ref)
	cmd.Dir = g.dbtProjectDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to checkout %s: %w", ref, err)
	}
	return nil
}

func (g *Operations) IsBehindOriginMain() (bool, error) {
	cmd := exec.Command("git", "rev-list", "--count", "HEAD..origin/main")
	cmd.Dir = g.dbtProjectDir
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to check if behind: %w", err)
	}

	count := strings.TrimSpace(string(output))
	return count != "0", nil
}

func (g *Operations) Rebase(ref string) error {
	cmd := exec.Command("git", "rebase", ref)
	cmd.Dir = g.dbtProjectDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to rebase onto %s: %w", ref, err)
	}
	return nil
}
