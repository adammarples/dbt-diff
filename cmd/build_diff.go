package cmd

import (
	"dbt-diff/internal/dbt"
	"dbt-diff/internal/git"
	"dbt-diff/internal/state"
	"fmt"
	"os"
)

// BuildDiff implements the build-diff command
func BuildDiff() error {
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

	fmt.Println("🔍 Analyzing changes...")

	// Get git information
	branch, err := gitOps.GetCurrentBranch()
	if err != nil {
		return err
	}

	diffHash, err := gitOps.GetDiffHash()
	if err != nil {
		return err
	}

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
		return fmt.Errorf("failed to stash changes: %w", err)
	}

	// Fetch origin
	fmt.Println("🌐 Fetching origin/main...")
	if err := gitOps.FetchOrigin(); err != nil {
		cleanup()
		return err
	}

	// Check if behind origin/main
	behind, err := gitOps.IsBehindOriginMain()
	if err != nil {
		cleanup()
		return err
	}

	if behind {
		fmt.Println("⚠️  Your branch is behind origin/main")
		fmt.Print("Would you like to rebase onto origin/main before continuing? (y/N): ")

		var response string
		fmt.Scanln(&response)

		if response == "y" || response == "Y" {
			fmt.Println("🔄 Rebasing onto origin/main...")
			if err := gitOps.Rebase("origin/main"); err != nil {
				cleanup()
				return fmt.Errorf("rebase failed: %w - please resolve conflicts and run again", err)
			}
			fmt.Println("✅ Rebase complete")
		}
	}

	// Get main SHA
	mainSha, err := gitOps.GetShortSHA()
	if err != nil {
		cleanup()
		return err
	}

	mainManifestPath := stateMgr.GetMainManifestPath(mainSha)

	// Compile main if needed
	if !stateMgr.ManifestExists(mainManifestPath) {
		fmt.Printf("📝 Compiling origin/main (%s)...\n", mainSha)

		if err := gitOps.Checkout("origin/main"); err != nil {
			cleanup()
			return err
		}

		if err := stateMgr.EnsureTargetDir(mainManifestPath); err != nil {
			cleanup()
			return err
		}

		if err := dbtRunner.Compile(mainManifestPath); err != nil {
			cleanup()
			stateMgr.RemovePartialManifest(mainManifestPath)
			return err
		}

		fmt.Println("✅ Main manifest compiled")
	} else {
		fmt.Printf("✅ Using cached main manifest (%s)\n", mainSha)
	}

	// Return to original branch
	fmt.Printf("🔄 Returning to branch %s...\n", branch)
	if err := gitOps.Checkout(branch); err != nil {
		cleanup()
		return err
	}

	// Apply stash
	fmt.Println("📤 Applying stashed changes...")
	if err := gitOps.PopStash(); err != nil {
		return fmt.Errorf("failed to apply stash: %w", err)
	}

	// Compile local
	localManifestPath := stateMgr.GetLocalManifestPath(diffHash)

	if !stateMgr.ManifestExists(localManifestPath) {
		fmt.Printf("📝 Compiling local changes (%s)...\n", diffHash[:8])

		if err := stateMgr.EnsureTargetDir(localManifestPath); err != nil {
			return err
		}

		if err := dbtRunner.Compile(localManifestPath); err != nil {
			stateMgr.RemovePartialManifest(localManifestPath)
			return err
		}

		fmt.Println("✅ Local manifest compiled")
	} else {
		fmt.Printf("✅ Using cached local manifest (%s)\n", diffHash[:8])
	}

	// Run modified models
	fmt.Println("🏗️  Running modified models...")
	if err := dbtRunner.Run(mainManifestPath); err != nil {
		return fmt.Errorf("run failed: %w", err)
	}

	fmt.Println("✅ Models run complete!")

	// Test modified models
	fmt.Println("🧪 Testing modified models...")
	if err := dbtRunner.Test(mainManifestPath); err != nil {
		return fmt.Errorf("test failed: %w", err)
	}

	fmt.Println("✅ Tests complete!")
	return nil
}
