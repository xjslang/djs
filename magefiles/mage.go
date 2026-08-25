//go:build mage

package main

import "github.com/magefile/mage/sh"

func Test() error {
	return sh.RunV("go", "test", "./...")
}

func Lint() error {
	return sh.RunV("golangci-lint", "run")
}

func LintFix() error {
	return sh.RunV("golangci-lint", "run", "--fix")
}

func Build() error {
	return sh.RunV("go", "build", "./cmd/djs")
}

func BuildWindows() error {
	return sh.RunWithV(map[string]string{
		"GOOS":   "windows",
		"GOARCH": "386",
	}, "go", "build", "./cmd/djs")
}

func BuildWasm() error {
	return sh.RunV("tinygo", "build", "-o", "djs.wasm", "-target", "wasm", "./cmd/wasm")
}

func Install() error {
	return sh.RunV("go", "install", "./cmd/djs/main.go")
}
