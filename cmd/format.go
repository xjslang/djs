package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/xjslang/djs/builder"
	"github.com/xjslang/djs/plugins"
	"github.com/xjslang/xjs/compiler"
	"github.com/xjslang/xjs/lexer"
)

func Format(args []string) int {
	flag.CommandLine.Parse(args)

	inputPath := flag.Arg(0)

	// Get original file permissions
	fileInfo, err := os.Stat(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", inputPath, err)
		return 1
	}
	originalMode := fileInfo.Mode()

	inputCode, err := os.ReadFile(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", inputPath, err)
		return 1
	}

	// Compile and transforms AST for proper formatting
	lb := lexer.NewBuilder()
	b := builder.New(lb).Build(string(inputCode))
	program, err := b.ParseProgram()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	program.Statements = plugins.Transform(program.Statements)
	result := compiler.New().WithPrettyPrint(compiler.WithSemi(false)).Compile(program)

	// Overwrite file preserving original permissions
	err = os.WriteFile(inputPath, []byte(result.Code), originalMode.Perm())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file %s: %v\n", inputPath, err)
		return 1
	}

	return 0
}
