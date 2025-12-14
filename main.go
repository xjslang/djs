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
	var outputPath string
	flag.StringVar(&outputPath, "o", "", "Output file path (transpile only, do not execute)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <file.djs>\n", filepath.Base(os.Args[0]))
		fmt.Fprintln(os.Stderr, "\nOptions:")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\nExamples:")
		fmt.Fprintln(os.Stderr, "  go run . test.djs           # Transpile and execute")
		fmt.Fprintln(os.Stderr, "  go run . test.djs -o out.js # Transpile to file")
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

	// Convert to absolute path for accurate source maps
	absInputPath, err := filepath.Abs(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving file path: %v\n", err)
		return 1
	}

	// Check if we're in transpile-only mode
	transpileOnly := outputPath != ""

	// Only check for Node if we're going to execute
	if !transpileOnly {
		if err := ensureNodeAvailable(); err != nil {
			fmt.Fprintf(os.Stderr, "Node.js not found: %v\n", err)
			return 1
		}
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

	// Only generate source map when executing (not when transpiling to file)
	c := compiler.New()
	if !transpileOnly {
		c = c.WithSourceMap()
	}
	result := c.Compile(program)

	// Transpile-only mode: write JS to output file (no source map)
	if transpileOnly {
		if err := ioutil.WriteFile(outputPath, []byte(result.Code), 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
			return 1
		}
		return 0
	}

	// Execute mode: prepare inline source map and run with Node
	// Enrich SourceMap with source metadata and file name
	sm := result.SourceMap
	if sm == nil {
		sm = &sourcemap.SourceMap{Version: 3}
	}
	// Set sources to the absolute path for proper error mapping
	sm.Sources = []string{absInputPath}
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
	jsBuilder.WriteString("//# sourceURL=" + absInputPath + "\n")
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
