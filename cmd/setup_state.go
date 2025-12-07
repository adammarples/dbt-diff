package cmd

import (
	"fmt"
	"os"

	"github.com/adammarples/dbt-diff/internal/dbt"
	"github.com/adammarples/dbt-diff/internal/git"
	"github.com/adammarples/dbt-diff/internal/state"
)

type StateInfo struct {
	MainManifestPath string
	MainSha          string
}

func SetupState(dbtOpts dbt.DbtOptions) (*StateInfo, error) {
	dbtProjectDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	stateMgr := state.New(dbtProjectDir)
	gitOps := git.New(dbtProjectDir)
	dbtRunner := dbt.NewWithOptions(dbtProjectDir, dbtOpts)

	if err := stateMgr.ValidateProjectRoot(); err != nil {
		return nil, err
	}
	if err := dbt.CheckDbtInstalled(); err != nil {
		return nil, err
	}

	fmt.Println("🌐 Fetching origin/main...")
	if err := gitOps.FetchOrigin(); err != nil {
		return nil, err
	}

	behind, err := gitOps.IsBehindOriginMain()
	if err != nil {
		return nil, err
	}

	if behind {
		fmt.Println("⚠️  Your branch is behind origin/main")
		fmt.Print("Would you like to rebase onto origin/main before continuing? (y/N): ")

		var response string
		_, _ = fmt.Scanln(&response)

		if response == "y" || response == "Y" {
			fmt.Println("🔄 Rebasing onto origin/main...")
			if err := gitOps.Rebase("origin/main"); err != nil {
				return nil, fmt.Errorf("rebase failed: %w - please resolve conflicts and run again", err)
			}
			fmt.Println("✅ Rebase complete")
		}
	}

	mainSha, err := gitOps.GetShortSHA()
	if err != nil {
		return nil, err
	}

	mainManifestPath := stateMgr.GetMainManifestPath(mainSha, dbtOpts.Target)

	// Only stash/checkout if we need to pull and compile main
	if !stateMgr.ManifestExists(mainManifestPath) {
		// Track whether we created a stash
		var stashCreated bool

		// Create cleanup function
		cleanup := func() {
			fmt.Println("🧹 Cleaning up...")
			if stashCreated {
				_ = gitOps.PopStash()
			}
			// Return to previous branch
			_ = gitOps.Checkout("-")
		}

		fmt.Println("📦 Stashing current changes...")
		if err := gitOps.CreateStash(); err != nil {
			return nil, fmt.Errorf("failed to stash changes: %w", err)
		}
		stashCreated = true

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
			_ = stateMgr.RemovePartialManifest(mainManifestPath)
			return nil, err
		}

		fmt.Println("✅ Main manifest compiled")

		// Return to previous branch
		fmt.Println("🔄 Returning to previous branch...")
		if err := gitOps.Checkout("-"); err != nil {
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

	return &StateInfo{
		MainManifestPath: mainManifestPath,
		MainSha:          mainSha,
	}, nil
}
