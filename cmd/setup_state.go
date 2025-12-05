package cmd

import (
	"fmt"
	"github.com/adammarples/dbt-diff/internal/dbt"
	"github.com/adammarples/dbt-diff/internal/git"
	"github.com/adammarples/dbt-diff/internal/state"
	"os"
)

// StateInfo contains the paths to compiled manifests
type StateInfo struct {
	MainManifestPath  string
	LocalManifestPath string
	MainSha           string
	DiffHash          string
}

// SetupState ensures both main and local manifests are compiled
// Returns state info needed for build/show commands
func SetupState() (*StateInfo, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	// Initialize components
	stateMgr := state.New(workDir)
	gitOps := git.New(workDir)
	dbtRunner := dbt.New(workDir)

	// Validate environment
	if err := stateMgr.ValidateProjectRoot(); err != nil {
		return nil, err
	}

	if err := dbt.CheckDbtInstalled(); err != nil {
		return nil, err
	}

	fmt.Println("🔍 Analyzing changes...")

	// Get git information
	branch, err := gitOps.GetCurrentBranch()
	if err != nil {
		return nil, err
	}

	diffHash, err := gitOps.GetDiffHash()
	if err != nil {
		return nil, err
	}

	// Fetch origin (doesn't require clean working directory)
	fmt.Println("🌐 Fetching origin/main...")
	if err := gitOps.FetchOrigin(); err != nil {
		return nil, err
	}

	// Check if behind origin/main
	behind, err := gitOps.IsBehindOriginMain()
	if err != nil {
		return nil, err
	}

	if behind {
		fmt.Println("⚠️  Your branch is behind origin/main")
		fmt.Print("Would you like to rebase onto origin/main before continuing? (y/N): ")

		var response string
		fmt.Scanln(&response)

		if response == "y" || response == "Y" {
			fmt.Println("🔄 Rebasing onto origin/main...")
			if err := gitOps.Rebase("origin/main"); err != nil {
				return nil, fmt.Errorf("rebase failed: %w - please resolve conflicts and run again", err)
			}
			fmt.Println("✅ Rebase complete")
		}
	}

	// Get main SHA (can get this without checking out)
	mainSha, err := gitOps.GetShortSHA()
	if err != nil {
		return nil, err
	}

	mainManifestPath := stateMgr.GetMainManifestPath(mainSha)

	// Only stash/checkout if we need to compile main
	if !stateMgr.ManifestExists(mainManifestPath) {
		stashName := gitOps.GetStashName(branch, diffHash)

		// Create cleanup function
		cleanup := func() {
			fmt.Println("🧹 Cleaning up...")
			exists, _ := gitOps.StashExists(stashName)
			if exists {
				_ = gitOps.PopStash()
			}
			// Return to original branch
			_ = gitOps.Checkout(branch)
		}

		// Create stash
		fmt.Println("📦 Stashing current changes...")
		if err := gitOps.CreateStash(stashName); err != nil {
			return nil, fmt.Errorf("failed to stash changes: %w", err)
		}

		fmt.Printf("📝 Compiling origin/main (%s)...\n", mainSha)

		if err := gitOps.Checkout("origin/main"); err != nil {
			cleanup()
			return nil, err
		}

		if err := stateMgr.EnsureTargetDir(mainManifestPath); err != nil {
			cleanup()
			return nil, err
		}

		if err := dbtRunner.Compile(mainManifestPath); err != nil {
			cleanup()
			stateMgr.RemovePartialManifest(mainManifestPath)
			return nil, err
		}

		fmt.Println("✅ Main manifest compiled")

		// Return to original branch
		fmt.Printf("🔄 Returning to branch %s...\n", branch)
		if err := gitOps.Checkout(branch); err != nil {
			cleanup()
			return nil, err
		}

		// Apply stash
		fmt.Println("📤 Applying stashed changes...")
		if err := gitOps.PopStash(); err != nil {
			return nil, fmt.Errorf("failed to apply stash: %w", err)
		}
	} else {
		fmt.Printf("✅ Using cached main manifest (%s)\n", mainSha)
	}

	// Compile local
	localManifestPath := stateMgr.GetLocalManifestPath(diffHash)

	if !stateMgr.ManifestExists(localManifestPath) {
		fmt.Printf("📝 Compiling local changes (%s)...\n", diffHash[:8])

		if err := stateMgr.EnsureTargetDir(localManifestPath); err != nil {
			return nil, err
		}

		if err := dbtRunner.Compile(localManifestPath); err != nil {
			stateMgr.RemovePartialManifest(localManifestPath)
			return nil, err
		}

		fmt.Println("✅ Local manifest compiled")
	} else {
		fmt.Printf("✅ Using cached local manifest (%s)\n", diffHash[:8])
	}

	return &StateInfo{
		MainManifestPath:  mainManifestPath,
		LocalManifestPath: localManifestPath,
		MainSha:           mainSha,
		DiffHash:          diffHash,
	}, nil
}
