package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	djsbuilder "github.com/xjslang/djs/builder"
	"github.com/xjslang/xjs/compiler"
	"github.com/xjslang/xjs/lexer"
)

const examplesDir = "../../examples"

// TestExamples verifies that all .djs files in the examples directory
// can be successfully parsed and transpiled without errors.
func TestExamples(t *testing.T) {
	// Find all .djs files in the examples directory
	djsFiles, err := findDJSFiles(examplesDir)
	if err != nil {
		t.Fatalf("Failed to find .djs files in examples directory: %v", err)
	}

	if len(djsFiles) == 0 {
		t.Fatal("No .djs files found in examples directory")
	}

	t.Logf("Found %d .djs files to verify", len(djsFiles))

	// Test each file
	for _, filePath := range djsFiles {
		t.Run(getTestName(filePath), func(t *testing.T) {
			// Read the file content
			content, err := os.ReadFile(filePath)
			if err != nil {
				t.Fatalf("Failed to read file %s: %v", filePath, err)
			}

			// Transpile the code
			jsCode, err := transpileDJSCode(string(content))
			if err != nil {
				t.Fatalf("Failed to transpile %s: %v", filePath, err)
			}

			// Verify that we got some output
			if strings.TrimSpace(jsCode) == "" {
				t.Errorf("Transpilation produced empty output for %s", filePath)
			}

			t.Logf("Successfully transpiled %s (%d bytes -> %d bytes)",
				filepath.Base(filePath), len(content), len(jsCode))
		})
	}
}

// findDJSFiles recursively finds all .djs files in the given directory
func findDJSFiles(root string) ([]string, error) {
	var djsFiles []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check if file has .djs extension
		if filepath.Ext(path) == ".djs" {
			djsFiles = append(djsFiles, path)
		}

		return nil
	})

	return djsFiles, err
}

// getTestName generates a readable test name from the file path
func getTestName(filePath string) string {
	// Get relative path from examples directory
	rel, err := filepath.Rel(examplesDir, filePath)
	if err != nil {
		rel = filepath.Base(filePath)
	}

	// Replace path separators with underscores and remove .djs extension
	name := strings.TrimSuffix(rel, ".djs")
	name = strings.ReplaceAll(name, string(filepath.Separator), "_")
	name = strings.ReplaceAll(name, "-", "_")

	return name
}

// transpileDJSCode transpiles DJS code to JavaScript
func transpileDJSCode(input string) (string, error) {
	lb := lexer.NewBuilder()
	p := djsbuilder.New(lb).Build(input)
	program, err := p.ParseProgram()
	if err != nil {
		return "", err
	}

	result := compiler.New().Compile(program)
	return result.Code, nil
}
