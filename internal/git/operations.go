package git

import (
	"crypto/sha256"
	"fmt"
	"os/exec"
	"strings"
)

// Operations handles git operations
type Operations struct {
	dbtProjectDir string
}

// New creates a new git operations handler
func New(dbtProjectDir string) *Operations {
	return &Operations{dbtProjectDir: dbtProjectDir}
}

// GetCurrentBranch returns the name of the current git branch
func (g *Operations) GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = g.dbtProjectDir
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// GetDiffHash computes a hash of the diff with origin/main
func (g *Operations) GetDiffHash() (string, error) {
	cmd := exec.Command("git", "diff", "origin/main")
	cmd.Dir = g.dbtProjectDir
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get diff: %w", err)
	}

	hash := sha256.Sum256(output)
	return fmt.Sprintf("%x", hash)[:8], nil
}

// GetShortSHA returns the short SHA of origin/main
func (g *Operations) GetShortSHA() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--short", "origin/main")
	cmd.Dir = g.dbtProjectDir
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get short SHA: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// CreateStash creates a stash with the given name
func (g *Operations) CreateStash(stashName string) error {
	cmd := exec.Command("git", "stash", "push", "-m", stashName)
	cmd.Dir = g.dbtProjectDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create stash: %w", err)
	}
	return nil
}

// PopStash pops the most recent stash
func (g *Operations) PopStash() error {
	cmd := exec.Command("git", "stash", "pop")
	cmd.Dir = g.dbtProjectDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to pop stash: %w", err)
	}
	return nil
}

// FetchOrigin fetches from origin
func (g *Operations) FetchOrigin() error {
	cmd := exec.Command("git", "fetch", "origin", "main")
	cmd.Dir = g.dbtProjectDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to fetch origin: %w", err)
	}
	return nil
}

// Checkout checks out a ref (branch or commit)
func (g *Operations) Checkout(ref string) error {
	cmd := exec.Command("git", "checkout", ref)
	cmd.Dir = g.dbtProjectDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to checkout %s: %w", ref, err)
	}
	return nil
}

// StashExists checks if a stash with the given name exists
func (g *Operations) StashExists(stashName string) (bool, error) {
	cmd := exec.Command("git", "stash", "list")
	cmd.Dir = g.dbtProjectDir
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to list stashes: %w", err)
	}
	return strings.Contains(string(output), stashName), nil
}

// GetStashName returns the formatted stash name
func (g *Operations) GetStashName(branch, diffHash string) string {
	return fmt.Sprintf("dbt-diff/%s/%s", branch, diffHash)
}

// IsBehindOriginMain checks if current branch is behind origin/main
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

// Rebase rebases current branch onto the given ref
func (g *Operations) Rebase(ref string) error {
	cmd := exec.Command("git", "rebase", ref)
	cmd.Dir = g.dbtProjectDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to rebase onto %s: %w", ref, err)
	}
	return nil
}
