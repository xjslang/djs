package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/xjslang/xjs/compiler"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/sourcemap"

	"github.com/xjslang/djs/plugins"
)

func main() {
	os.Exit(run())
}

func run() int {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <file.djs>\n", filepath.Base(os.Args[0]))
		fmt.Fprintln(os.Stderr, "Transpiles DJS to JS and executes with Node.")
		fmt.Fprintln(os.Stderr, "Example: go run . test.djs")
	}

	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
		return 2
	}

	inputPath := flag.Arg(0)
	inputCode, err := ioutil.ReadFile(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		return 1
	}

	if err := ensureNodeAvailable(); err != nil {
		fmt.Fprintf(os.Stderr, "Node.js not found: %v\n", err)
		return 1
	}

	lb := lexer.NewBuilder()
	p := parser.NewBuilder(lb).
		Install(plugins.DeferPlugin).
		Install(plugins.OrPlugin).
		Install(plugins.StrictEqualityPlugin).
		Build(string(inputCode))

	program, perr := p.ParseProgram()
	if perr != nil {
		fmt.Fprintln(os.Stderr, perr)
		return 1
	}

	result := compiler.New().
		WithSourceMap().
		Compile(program)

	// Enrich SourceMap with source metadata and file name
	sm := result.SourceMap
	if sm == nil {
		sm = &sourcemap.SourceMap{Version: 3}
	}
	// Set sources to the original file path and inline content
	sm.Sources = []string{inputPath}
	sm.SourcesContent = []string{string(inputCode)}

	// Determine output file name (for tooling); not strictly needed for inline maps
	outFile := deriveOutputFilename(inputPath)
	sm.File = filepath.Base(outFile)

	// Serialize SourceMap to base64 JSON and embed as inline comment
	smJSON, jerr := json.Marshal(sm)
	if jerr != nil {
		fmt.Fprintf(os.Stderr, "Error serializing source map: %v\n", jerr)
		return 1
	}
	b64 := base64.StdEncoding.EncodeToString(smJSON)

	// Compose final JS with inline source map and optional sourceURL
	var jsBuilder strings.Builder
	jsBuilder.WriteString(result.Code)
	if !strings.HasSuffix(result.Code, "\n") {
		jsBuilder.WriteString("\n")
	}
	// Help debuggers: set a sourceURL for nicer stack display
	jsBuilder.WriteString("//# sourceURL=" + inputPath + "\n")
	jsBuilder.WriteString("//# sourceMappingURL=data:application/json;charset=utf-8;base64,")
	jsBuilder.WriteString(b64)
	jsBuilder.WriteString("\n")

	finalJS := jsBuilder.String()

	// Write to a temporary JS file so Node can execute it with source maps
	tmpFile, terr := writeTempJS(outFile, finalJS)
	if terr != nil {
		fmt.Fprintf(os.Stderr, "Error writing temp JS: %v\n", terr)
		return 1
	}
	defer os.Remove(tmpFile)

	// Execute with Node enabling source maps so runtime errors map to original DJS
	cmd := exec.Command("node", "--enable-source-maps", tmpFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		// Preserve Node’s exit code when possible
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			if exitErr.ProcessState != nil {
				return exitErr.ProcessState.ExitCode()
			}
			return 1
		}
		fmt.Fprintf(os.Stderr, "Execution error: %v\n", err)
		return 1
	}

	return 0
}

func ensureNodeAvailable() error {
	cmd := exec.Command("node", "--version")
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}

func deriveOutputFilename(inputPath string) string {
	base := filepath.Base(inputPath)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)
	return name + ".transpiled.js"
}

func writeTempJS(outFileName, content string) (string, error) {
	// Use the OS temp dir but provide a stable-ish name for readability
	tmpDir := os.TempDir()
	tmpPath := filepath.Join(tmpDir, outFileName)
	if err := ioutil.WriteFile(tmpPath, []byte(content), 0o644); err != nil {
		return "", err
	}
	return tmpPath, nil
}
