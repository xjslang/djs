//go:build mage

package main

import "github.com/magefile/mage/sh"

func Lint() error {
	return sh.RunV("golangci-lint", "run")
}

func LintFix() error {
	return sh.RunV("golangci-lint", "run", "--fix")
}

func BuildWasm() error {
	return sh.RunV("tinygo", "build", "-o", "djs.wasm", "-target", "wasm", "./cmd/wasm")
}
