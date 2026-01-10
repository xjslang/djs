//go:build mage

package main

import (
	"fmt"
	"time"

	"github.com/magefile/mage/sh"
)

// Test runs all Go tests in the project using 'go test ./...'
func Test() error {
	fmt.Println("🧪 Running tests...")
	return sh.RunV("go", "test", "./...")
}

// Bench runs project benchmarks using 'go test -bench=.'
func Bench() error {
	fmt.Println("⚡ Running benchmarks...")
	return sh.RunV("go", "test", "-bench=.")
}

// Tidy cleans and organizes the go.mod file using 'go mod tidy'
func Tidy() error {
	fmt.Println("🔧 Tidying go.mod...")
	return sh.RunV("go", "mod", "tidy")
}

// Lint runs linting
func Lint() error {
	fmt.Println("🔍 Running linter...")
	return sh.RunV("golangci-lint", "run")
}

// InstallHooks configures Git hooks for the project
func InstallHooks() error {
	fmt.Println("🔗 Installing Git hooks...")
	if err := sh.RunV("git", "config", "core.hooksPath", ".githooks"); err != nil {
		return fmt.Errorf("failed to set hooks path: %w", err)
	}
	if err := sh.RunV("chmod", "+x", ".githooks/pre-push"); err != nil {
		return fmt.Errorf("failed to make pre-push executable: %w", err)
	}
	fmt.Println("✅ Git hooks installed successfully!")
	return nil
}

// Clean removes temporary files and cache
func Clean() error {
	fmt.Println("🧹 Cleaning temporary files and cache...")
	if err := sh.RunV("go", "clean", "-testcache"); err != nil {
		fmt.Println("Note: failed to clean test cache, continuing...")
	}
	if err := sh.RunV("go", "clean", "-modcache"); err != nil {
		fmt.Println("Note: failed to clean mod cache, continuing...")
	}
	fmt.Println("✅ Cleanup completed!")
	return nil
}

// Install installs dependencies
func Install() error {
	fmt.Println("📦 Installing dependencies...")
	return sh.RunV("go", "mod", "download")
}

// Build compiles the project with version information
func Build() error {
	fmt.Println("🔨 Building djs...")

	// Get version from git tag
	version, err := sh.Output("git", "describe", "--tags", "--always", "--dirty")
	if err != nil {
		version = "dev"
	}

	// Get git commit hash
	commit, err := sh.Output("git", "rev-parse", "--short", "HEAD")
	if err != nil {
		commit = "unknown"
	}

	// Get build date
	buildDate := time.Now().UTC().Format("2006-01-02_15:04:05")

	// Build with ldflags
	ldflags := fmt.Sprintf(
		"-X main.Version=%s -X main.GitCommit=%s -X main.BuildDate=%s",
		version, commit, buildDate,
	)

	return sh.RunV("go", "build", "-ldflags", ldflags, "-o", "djs", ".")
}
