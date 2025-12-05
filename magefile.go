//go:build mage

package main

import (
	"github.com/magefile/mage/sh"
)

// Bench runs project benchmarks using 'go test -bench=.'
func Bench() error {
	return sh.RunV("go", "test", "-bench=.")
}

// Lint runs linting
func Lint() error {
	return sh.RunV("golangci-lint", "run")
}
