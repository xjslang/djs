package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/xjslang/djs/builder"
	"github.com/xjslang/djs/plugins"
	"github.com/xjslang/xjs/compiler"
	"github.com/xjslang/xjs/lexer"
)

func Format(args []string) int {
	var stdin bool
	flag.BoolVar(&stdin, "stdin", false, "Read input from stdin, instead of a file, and prints it formatted")
	_ = flag.CommandLine.Parse(args)

	inputPath := flag.Arg(0)

	if stdin {
		// Read content from stdin
		reader := bufio.NewReader(os.Stdin)
		var builder strings.Builder
		_, err := io.Copy(&builder, reader)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
			return 1
		}
		data := builder.String()

		// Compile and transforms AST for proper formatting
		output, err := compile(data)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}

		fmt.Println(output)
	} else {
		// Read file contents
		data, err := os.ReadFile(inputPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", inputPath, err)
			return 1
		}

		// Compile and transforms AST for proper formatting
		output, err := compile(string(data))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}

		// Save file and preserve the original permissions
		err = saveFile(inputPath, []byte(output))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file %s: %v\n", inputPath, err)
			return 1
		}
	}

	return 0
}

// Compile and transforms AST for proper formatting
func compile(input string) (string, error) {
	lb := lexer.NewBuilder()
	b := builder.New(lb).Build(input)
	program, err := b.ParseProgram()
	if err != nil {
		return "", err
	}
	program.Statements = plugins.Transform(program.Statements)
	result := compiler.New().WithPrettyPrint(compiler.WithSemi(false)).Compile(program)
	return result.Code, nil
}

// Save file contents and preserve original permissions
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
