//go:build mage

package main

import (
	"fmt"

	"github.com/magefile/mage/sh"
)

// Test runs all Go tests in the project using 'go test ./...'
func Test() error {
	return sh.RunV("go", "test", "./...")
}

// Bench runs project benchmarks using 'go test -bench=.'
func Bench() error {
	return sh.RunV("go", "test", "-bench=.")
}

// Tidy cleans and organizes the go.mod file using 'go mod tidy'
func Tidy() error {
	return sh.RunV("go", "mod", "tidy")
}

// Lint runs linting
func Lint() error {
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
