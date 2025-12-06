package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/xjslang/djs/plugins"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
)

func execCommand(source string) error {
	cmd := exec.Command("node", "-e", source)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func main() {
	filename := os.Args[1]
	input, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	lb := lexer.NewBuilder()
	p := parser.NewBuilder(lb).
		Install(plugins.DeferPlugin).
		Install(plugins.OrPlugin).
		Build(string(input))
	program, err := p.ParseProgram()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	execCommand(program.String())
}
