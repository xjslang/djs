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
	var generateSourceMap bool
	var inlineSourceMap bool
	var inlineSources bool
	var mapRoot string
	var sourceRoot string
	flag.StringVar(&outputPath, "o", "", "Output file path (transpile only, do not execute)")
	flag.BoolVar(&generateSourceMap, "sourcemap", false, "Generate external source map file (.map)")
	flag.BoolVar(&inlineSourceMap, "inline-sourcemap", false, "Embed source map as base64 in output file")
	flag.BoolVar(&inlineSources, "inline-sources", false, "Include source content in source map")
	flag.StringVar(&mapRoot, "map-root", "", "Root path for source map file location")
	flag.StringVar(&sourceRoot, "source-root", "", "Root path for source files (sourceRoot field in map)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <file.djs>\n", filepath.Base(os.Args[0]))
		fmt.Fprintln(os.Stderr, "\nOptions:")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\nExamples:")
		fmt.Fprintln(os.Stderr, "  go run . test.djs                                           # Transpile and execute")
		fmt.Fprintln(os.Stderr, "  go run . -o out.js test.djs                                 # Transpile to file")
		fmt.Fprintln(os.Stderr, "  go run . -o out.js --sourcemap test.djs                     # External source map")
		fmt.Fprintln(os.Stderr, "  go run . -o out.js --inline-sourcemap test.djs              # Embedded source map")
		fmt.Fprintln(os.Stderr, "  go run . -o out.js --sourcemap --inline-sources test.djs    # With embedded sources")
		fmt.Fprintln(os.Stderr, "  go run . -o out.js --sourcemap --map-root /maps/ test.djs   # Map in /maps/ folder")
		fmt.Fprintln(os.Stderr, "  go run . -o out.js --sourcemap --source-root /src/ test.djs # Source root prefix")
	}

	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
		return 2
	}

	// Validate mutually exclusive flags
	if generateSourceMap && inlineSourceMap {
		fmt.Fprintln(os.Stderr, "Error: --sourcemap and --inline-sourcemap are mutually exclusive")
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

	// Generate source map when executing OR when explicitly requested
	c := compiler.New()
	if !transpileOnly || generateSourceMap || inlineSourceMap {
		c = c.WithSourceMap()
	}
	result := c.Compile(program)

	// Transpile-only mode: write JS to output file
	if transpileOnly {
		if generateSourceMap || inlineSourceMap {
			// Prepare source map
			sm := result.SourceMap
			if sm == nil {
				sm = &sourcemap.SourceMap{Version: 3}
			}
			sm.Sources = []string{absInputPath}
			if inlineSources {
				sm.SourcesContent = []string{string(inputCode)}
			}
			// Set sourceRoot field in source map JSON
			if sourceRoot != "" {
				sm.SourceRoot = sourceRoot
			}
			sm.File = filepath.Base(outputPath)

			if generateSourceMap {
				// External source map file
				mapPath := outputPath + ".map"
				smJSON, jerr := json.MarshalIndent(sm, "", "  ")
				if jerr != nil {
					fmt.Fprintf(os.Stderr, "Error serializing source map: %v\n", jerr)
					return 1
				}
				if err := ioutil.WriteFile(mapPath, smJSON, 0o644); err != nil {
					fmt.Fprintf(os.Stderr, "Error writing source map file: %v\n", err)
					return 1
				}

				// Write JS file with external source map reference
				var jsBuilder strings.Builder
				jsBuilder.WriteString(result.Code)
				if !strings.HasSuffix(result.Code, "\n") {
					jsBuilder.WriteString("\n")
				}
				jsBuilder.WriteString("//# sourceMappingURL=")
				// Apply map root to the URL of the .map file
				if mapRoot != "" {
					jsBuilder.WriteString(mapRoot)
					if !strings.HasSuffix(mapRoot, "/") {
						jsBuilder.WriteString("/")
					}
				}
				jsBuilder.WriteString(filepath.Base(mapPath))
				jsBuilder.WriteString("\n")

				if err := ioutil.WriteFile(outputPath, []byte(jsBuilder.String()), 0o644); err != nil {
					fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
					return 1
				}
			} else if inlineSourceMap {
				// Inline source map (base64 embedded)
				smJSON, jerr := json.Marshal(sm)
				if jerr != nil {
					fmt.Fprintf(os.Stderr, "Error serializing source map: %v\n", jerr)
					return 1
				}
				b64 := base64.StdEncoding.EncodeToString(smJSON)

				var jsBuilder strings.Builder
				jsBuilder.WriteString(result.Code)
				if !strings.HasSuffix(result.Code, "\n") {
					jsBuilder.WriteString("\n")
				}
				jsBuilder.WriteString("//# sourceMappingURL=data:application/json;charset=utf-8;base64,")
				jsBuilder.WriteString(b64)
				jsBuilder.WriteString("\n")

				if err := ioutil.WriteFile(outputPath, []byte(jsBuilder.String()), 0o644); err != nil {
					fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
					return 1
				}
			}
		} else {
			// Write JS file without source map
			if err := ioutil.WriteFile(outputPath, []byte(result.Code), 0o644); err != nil {
				fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
				return 1
			}
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
	// Always include sources in execution mode for better error messages
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
