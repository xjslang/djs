package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dop251/goja"
	"github.com/xjslang/xjs/compiler"
	"github.com/xjslang/xjs/lexer"

	djsbuilder "github.com/xjslang/djs/builder"
)

const testDataDir = "../../testdata"

type TranspilationTest struct {
	name           string
	inputFile      string
	expectedOutput string
}

func normalizeLineEndings(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	return s
}

func transpileXJSCode(input string) (string, error) {
	lb := lexer.NewBuilder()
	p := djsbuilder.New(lb).Build(input)
	program, err := p.ParseProgram()
	if err != nil {
		return "", fmt.Errorf("ParseProgram error: %v", err)
	}

	// Convert the AST to JavaScript code (now with automatic semicolons)
	result := compiler.New().Compile(program)
	return result.Code, nil
}

func executeJavaScript(code string) (string, error) {
	vm := goja.New()
	var output strings.Builder
	_ = vm.Set("console", map[string]any{
		"log": func(args ...any) {
			for i, arg := range args {
				if i > 0 {
					output.WriteString(" ")
				}
				if arg == nil {
					output.WriteString("null")
				} else {
					output.WriteString(fmt.Sprintf("%v", arg))
				}
			}
			output.WriteString("\n")
		},
	})
	_, err := vm.RunString(code)
	if err != nil {
		return "", fmt.Errorf("failed to execute JavaScript: %v", err)
	}
	result := strings.TrimSpace(output.String())
	return normalizeLineEndings(result), nil
}

func loadTestCase(t *testing.T, baseName string) TranspilationTest {
	inputFile := filepath.Join(testDataDir, baseName+".djs")
	outputFile := filepath.Join(testDataDir, baseName+".output")

	// Read input file
	inputContent, err := os.ReadFile(inputFile)
	if err != nil {
		t.Fatalf("Failed to read input file %s: %v", inputFile, err)
	}

	// Read expected output file
	expectedOutput, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file %s: %v", outputFile, err)
	}

	return TranspilationTest{
		name:           baseName,
		inputFile:      string(inputContent),
		expectedOutput: normalizeLineEndings(strings.TrimSpace(string(expectedOutput))),
	}
}

func RunTranspilationTest(t *testing.T, test TranspilationTest) {
	t.Run(test.name, func(t *testing.T) {
		// Transpile the XJS code to JavaScript
		transpiledJS, err := transpileXJSCode(test.inputFile)
		if err != nil {
			t.Fatalf("Transpilation failed: %v", err)
		}

		// Execute the transpiled JavaScript
		actualOutput, err := executeJavaScript(transpiledJS)
		if err != nil {
			t.Fatalf("JavaScript execution failed: %v", err)
		}

		// Compare the actual output with expected output
		actualOutput = normalizeLineEndings(strings.TrimSpace(actualOutput))
		if actualOutput != test.expectedOutput {
			t.Errorf("Output mismatch:\nExpected: %q\nActual:   %q\nTranspiled JS:\n%s",
				test.expectedOutput, actualOutput, transpiledJS)
		}
	})
}

func TestTranspilation(t *testing.T) {
	// discover test cases by recursively
	var testCases []string
	err := filepath.WalkDir(testDataDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(d.Name(), ".djs") {
			// Get relative path from testDataDir
			relPath, err := filepath.Rel(testDataDir, path)
			if err != nil {
				return err
			}

			// Remove .djs extension to get the test case name
			testCaseName := strings.TrimSuffix(relPath, ".djs")

			// Check if corresponding .output file exists
			outputFile := filepath.Join(testDataDir, testCaseName+".output")
			if _, err := os.Stat(outputFile); err == nil {
				testCases = append(testCases, testCaseName)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Failed to walk testdata directory: %v", err)
	}

	if len(testCases) == 0 {
		t.Fatal("No test cases found in testdata directory")
	}

	for _, testCase := range testCases {
		test := loadTestCase(t, testCase)
		RunTranspilationTest(t, test)
	}
}

// TestAsyncAwaitErrors tests error handling for async/await
func TestAsyncAwaitErrors(t *testing.T) {
	errorTests := []struct {
		name  string
		input string
	}{
		{
			name:  "await_at_top_level",
			input: "let x = await get()",
		},
		{
			name:  "await_at_top_level_in_expression",
			input: "console.log(await fetch('data'))",
		},
		{
			name: "multiple_await_at_top_level",
			input: `
let x = await fetch('https://api.example.com')
let y = await fetch('https://api2.example.com')
console.log(x, y)
`,
		},
	}

	for _, test := range errorTests {
		t.Run(test.name, func(t *testing.T) {
			_, err := transpileXJSCode(test.input)
			if err == nil {
				t.Errorf("Expected transpilation error for input: %s", test.input)
			}
		})
	}
}
