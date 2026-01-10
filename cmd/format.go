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

	// Read file contents
	data, err := os.ReadFile(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", inputPath, err)
		return 1
	}

	// Compile and transforms AST for proper formatting
	lb := lexer.NewBuilder()
	b := builder.New(lb).Build(string(data))
	program, err := b.ParseProgram()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	program.Statements = plugins.Transform(program.Statements)
	result := compiler.New().WithPrettyPrint(compiler.WithSemi(false)).Compile(program)

	// Save file preserving the original permissions
	err = saveFile(inputPath, []byte(result.Code))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file %s: %v\n", inputPath, err)
		return 1
	}

	return 0
}

func saveFile(path string, data []byte) error {
	// Get original file permissions
	fileInfo, err := os.Stat(path)
	if err != nil {
		return err
	}
	originalMode := fileInfo.Mode()

	// Overwrite file preserving original permissions
	err = os.WriteFile(path, data, originalMode.Perm())
	if err != nil {
		return err
	}

	return nil
}
