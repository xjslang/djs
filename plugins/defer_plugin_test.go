package plugins

import (
	"fmt"
	"testing"

	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
)

func TestDeferOutsideFunction(t *testing.T) {
	input := `
	defer {
		console.log("This should cause an error")
	}`
	lb := lexer.NewBuilder()
	p := parser.NewBuilder(lb).Install(DeferPlugin).Build(input)
	_, err := p.ParseProgram() // Parse to trigger error checking
	if err == nil {
		t.Errorf("Expected error when defer is used outside function, but got none")
	}
}

func TestDeferInsideNestedFunction(t *testing.T) {
	input := `
	function outer() {
		function inner() {
			defer {
				console.log("This should work")
			}
		}
	}`
	lb := lexer.NewBuilder()
	p := parser.NewBuilder(lb).Install(DeferPlugin).Build(input)
	_, err := p.ParseProgram() // Parse to trigger error checking
	if err != nil {
		t.Errorf("Expected error when defer is used outside function, but got none")
	}

	fmt.Println("Nested function defer parsed successfully")
}

func TestDeferSemicolonHandling(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name: "explicit semicolons",
			input: `
	function cleanup() {
		defer closeFile();
		defer cleanupResources();
	}`,
		},
		{
			name: "ASI (Automatic Semicolon Insertion)",
			input: `
	function cleanup() {
		defer closeFile()
		defer cleanupResources()
	}`,
		},
		{
			name: "block without semicolon",
			input: `
	function cleanup() {
		defer {
			closeFile()
			cleanupResources()
		}
		let x = 1
	}`,
		},
		{
			name: "mixed defer styles",
			input: `
	function cleanup() {
		defer closeFile();
		defer {
			cleanupA()
			cleanupB()
		}
		defer cleanupC()
	}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lb := lexer.NewBuilder()
			p := parser.NewBuilder(lb).Install(DeferPlugin).Build(tt.input)
			prog, err := p.ParseProgram()
			if err != nil {
				t.Fatalf("Expected no error for %s, got: %v", tt.name, err)
			}
			if prog == nil {
				t.Fatalf("Expected program to be parsed for %s", tt.name)
			}
		})
	}
}
