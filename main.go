package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rs/xid"
	"github.com/xjslang/xjs/compiler"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/sourcemap"

	djsbuilder "github.com/xjslang/djs/builder"
)

type ParserErrors struct {
	Errors []parser.ParserError `json:"errors"`
}

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
	var jsonOutput bool
	var checkOnly bool
	var useStdin bool
	var stdinFilename string
	flag.StringVar(&outputPath, "o", "", "Output file path (transpile only, do not execute)")
	flag.BoolVar(&generateSourceMap, "sourcemap", false, "Generate external source map file (.map)")
	flag.BoolVar(&inlineSourceMap, "inline-sourcemap", false, "Embed source map as base64 in output file")
	flag.BoolVar(&inlineSources, "inline-sources", false, "Include source content in source map")
	flag.StringVar(&mapRoot, "map-root", "", "Root path for source map file location")
	flag.StringVar(&sourceRoot, "source-root", "", "Root path for source files (sourceRoot field in map)")
	flag.BoolVar(&jsonOutput, "json", false, "Output errors in JSON format")
	flag.BoolVar(&checkOnly, "check", false, "Check syntax only, do not execute or transpile")
	flag.BoolVar(&useStdin, "stdin", false, "Read source code from stdin instead of file")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <file.djs>\n", filepath.Base(os.Args[0]))
		fmt.Fprintln(os.Stderr, "\nOptions:")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\nExamples:")
		fmt.Fprintln(os.Stderr, "  djs input.djs                                                   # Transpile and execute")
		fmt.Fprintln(os.Stderr, "  djs --check input.djs                                           # Check syntax only")
		fmt.Fprintln(os.Stderr, "  djs --check --json input.djs                                    # Check syntax, output JSON")
		fmt.Fprintln(os.Stderr, "  cat input.djs | djs --stdin --check --json                      # Check from stdin")
		fmt.Fprintln(os.Stderr, "  djs --json input.djs                                            # Show errors in JSON format")
		fmt.Fprintln(os.Stderr, "  djs -o output.js input.djs                                      # Transpile to file")
		fmt.Fprintln(os.Stderr, "  djs -o output.js --sourcemap input.djs                          # External source map")
		fmt.Fprintln(os.Stderr, "  djs -o output.js --inline-sourcemap input.djs                   # Embedded source map")
		fmt.Fprintln(os.Stderr, "  djs -o output.js --sourcemap --inline-sources input.djs         # With embedded sources")
		fmt.Fprintln(os.Stderr, "  djs -o output.js --sourcemap --map-root /maps/ input.djs        # Map in /maps/ folder")
		fmt.Fprintln(os.Stderr, "  djs -o output.js --sourcemap --source-root /src/ input.djs      # Source root prefix")
	}

	flag.Parse()

	// Validate stdin usage
	if useStdin {
		if flag.NArg() != 0 {
			fmt.Fprintln(os.Stderr, "Error: --stdin cannot be used with a file argument")
			return 2
		}
	} else {
		if flag.NArg() != 1 {
			flag.Usage()
			return 2
		}
	}

	// Validate mutually exclusive flags
	if generateSourceMap && inlineSourceMap {
		fmt.Fprintln(os.Stderr, "Error: --sourcemap and --inline-sourcemap are mutually exclusive")
		return 2
	}

	// Validate --check is incompatible with -o and source map flags
	if checkOnly && outputPath != "" {
		fmt.Fprintln(os.Stderr, "Error: --check cannot be used with -o (transpile mode)")
		return 2
	}
	if checkOnly && (generateSourceMap || inlineSourceMap || inlineSources || mapRoot != "" || sourceRoot != "") {
		fmt.Fprintln(os.Stderr, "Error: --check cannot be used with source map flags")
		return 2
	}

	// Check if we're in transpile-only mode
	transpileOnly := outputPath != ""
	hasSourceMap := generateSourceMap || inlineSourceMap

	// Validate --map-root requires --sourcemap (not --inline-sourcemap)
	if mapRoot != "" && !generateSourceMap {
		fmt.Fprintln(os.Stderr, "Error: --map-root requires --sourcemap (external source map)")
		return 2
	}

	// Validate --inline-sources requires some form of source map
	if inlineSources && !hasSourceMap {
		fmt.Fprintln(os.Stderr, "Error: --inline-sources requires --sourcemap or --inline-sourcemap")
		return 2
	}

	// Validate --source-root requires some form of source map
	if sourceRoot != "" && !hasSourceMap {
		fmt.Fprintln(os.Stderr, "Error: --source-root requires --sourcemap or --inline-sourcemap")
		return 2
	}

	// Validate source map flags require transpile mode (-o flag)
	if !transpileOnly && (generateSourceMap || inlineSourceMap || inlineSources || mapRoot != "" || sourceRoot != "") {
		fmt.Fprintln(os.Stderr, "Error: source map flags require -o (transpile mode)")
		return 2
	}

	var inputCode []byte
	var absInputPath string
	var err error

	if useStdin {
		// Read from stdin
		reader := bufio.NewReader(os.Stdin)
		var builder strings.Builder
		_, err = io.Copy(&builder, reader)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
			return 1
		}
		inputCode = []byte(builder.String())
		absInputPath = stdinFilename
	} else {
		// Read from file
		inputPath := flag.Arg(0)
		inputCode, err = os.ReadFile(inputPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			return 1
		}

		// Convert to absolute path for accurate source maps
		absInputPath, err = filepath.Abs(inputPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error resolving file path: %v\n", err)
			return 1
		}
	}

	// Only check for Node if we're going to execute
	if !transpileOnly && !checkOnly {
		if err := ensureNodeAvailable(); err != nil {
			fmt.Fprintf(os.Stderr, "Node.js not found: %v\n", err)
			return 1
		}
	}

	lb := lexer.NewBuilder()
	p := djsbuilder.New(lb).Build(string(inputCode))

	program, perr := p.ParseProgram()
	if perr != nil {
		if jsonOutput {
			// Output errors in JSON format using parser.ParserErrors
			parserErrs := ParserErrors{
				Errors: p.Errors(),
			}

			jsonBytes, jerr := json.MarshalIndent(parserErrs, "", "  ")
			if jerr != nil {
				fmt.Fprintf(os.Stderr, "Error serializing error response: %v\n", jerr)
				return 1
			}
			fmt.Fprintln(os.Stdout, string(jsonBytes))
		} else {
			fmt.Fprintln(os.Stderr, perr)
		}
		return 1
	}

	// Check-only mode: if no errors, exit successfully
	if checkOnly {
		return 0
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
				if err := os.WriteFile(mapPath, smJSON, 0o644); err != nil {
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

				if err := os.WriteFile(outputPath, []byte(jsBuilder.String()), 0o644); err != nil {
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

				if err := os.WriteFile(outputPath, []byte(jsBuilder.String()), 0o644); err != nil {
					fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
					return 1
				}
			}
		} else {
			// Write JS file without source map
			if err := os.WriteFile(outputPath, []byte(result.Code), 0o644); err != nil {
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
	outFile := deriveOutputFilename(absInputPath)
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
	// Use the directory of the original DJS file to preserve require() resolution
	inputDir := filepath.Dir(absInputPath)
	tmpFile, terr := writeTempJS(inputDir, outFile, finalJS)
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
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("node command not found or failed to execute")
	}

	version := strings.TrimSpace(string(output))
	// Remove 'v' prefix if present (e.g., "v18.16.0" -> "18.16.0")
	version = strings.TrimPrefix(version, "v")

	// Parse version (format: major.minor.patch)
	var major, minor int
	_, err = fmt.Sscanf(version, "%d.%d", &major, &minor)
	if err != nil {
		return fmt.Errorf("unable to parse node version: %s", version)
	}

	// DJS requires Node.js 7.6+ for async/await support
	if major < 7 || (major == 7 && minor < 6) {
		return fmt.Errorf("node version %s is too old; DJS requires Node.js 7.6 or later for async/await support", version)
	}

	return nil
}

func deriveOutputFilename(inputPath string) string {
	base := filepath.Base(inputPath)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)
	return name + ".transpiled.js"
}

func writeTempJS(baseDir, outFileName, content string) (string, error) {
	// Write temp file in the same directory as the original DJS file
	// so Node.js can resolve relative requires and local node_modules
	ext := filepath.Ext(outFileName)
	base := strings.TrimSuffix(outFileName, ext)
	uniqueFileName := base + "." + xid.New().String() + ext
	tmpPath := filepath.Join(baseDir, uniqueFileName)
	if err := os.WriteFile(tmpPath, []byte(content), 0o644); err != nil {
		return "", err
	}
	return tmpPath, nil
}
