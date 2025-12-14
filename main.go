package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/xjslang/djs/plugins"
	"github.com/xjslang/xjs/compiler"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
)

func execCommand(source string) error {
	cmd := exec.Command("node", "--enable-source-maps", source)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func main() {
	inputFile := os.Args[1]
	bytes, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	input := string(bytes)

	lb := lexer.NewBuilder()
	p := parser.NewBuilder(lb).
		Install(plugins.DeferPlugin).
		Install(plugins.OrPlugin).
		Install(plugins.StrictEqualityPlugin).
		Build(input)

	program, err := p.ParseProgram()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	result := compiler.New().
		WithSourceMap().
		Compile(program)

	// Prepare output file paths
	jsOutFile := inputFile + ".js"
	mapOutFile := inputFile + ".map"

	// Add source map reference to the JS output and write files
	jsCode := result.Code + "\n//# sourceMappingURL=" + mapOutFile
	sourceMap := result.SourceMap
	sourceMap.Sources = []string{inputFile}
	mapJSON, _ := json.Marshal(sourceMap)

	// Write source map
	err = os.WriteFile(mapOutFile, mapJSON, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Write compiled JS
	err = os.WriteFile(jsOutFile, []byte(jsCode), 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Execute Node.js on the output file
	err = execCommand(jsOutFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
