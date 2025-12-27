package plugins

import (
	"testing"

	"github.com/xjslang/xjs/compiler"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
)

func TestThrowWithArgument(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "throw string literal",
			input:    `throw "error message"`,
			expected: `throw "error message";`,
		},
		{
			name:     "throw variable",
			input:    `throw err`,
			expected: `throw err;`,
		},
		{
			name:     "throw expression",
			input:    `throw "Error: " + message`,
			expected: `throw ("Error: "+message);`,
		},
		{
			name: "throw in function",
			input: `function test() {
				throw "error"
			}`,
			expected: `function test(){throw "error";}`,
		},
		{
			name:     "throw with new expression",
			input:    `throw new Error("message")`,
			expected: `throw new Error("message");`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lb := lexer.NewBuilder()
			p := parser.NewBuilder(lb).
				Install(ThrowPlugin).
				Install(NewPlugin).
				Build(tt.input)
			prog, err := p.ParseProgram()
			if err != nil {
				t.Fatalf("Expected no error, got: %v", err)
			}
			if prog == nil {
				t.Fatal("Expected program to be parsed")
			}

			result := compiler.New().Compile(prog)
			if result.Code != tt.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tt.expected, result.Code)
			}
		})
	}
}

func TestThrowWithoutArgument(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "throw with semicolon",
			input: `throw;`,
		},
		{
			name:  "throw at end of block",
			input: `function test() { throw }`,
		},
		{
			name:  "throw at end of file",
			input: `throw`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lb := lexer.NewBuilder()
			p := parser.NewBuilder(lb).
				Install(ThrowPlugin).
				Build(tt.input)
			_, err := p.ParseProgram()
			if err == nil {
				t.Error("Expected error when throw has no argument, but got none")
			}
		})
	}
}

func TestThrowWithASI(t *testing.T) {
	// throw with newline (ASI)
	input := `
	function test() {
		throw "error"
		let x = 1
	}`
	lb := lexer.NewBuilder()
	p := parser.NewBuilder(lb).
		Install(ThrowPlugin).
		Build(input)
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("Expected no error with ASI, got: %v", err)
	}
	if prog == nil {
		t.Fatal("Expected program to be parsed")
	}
}
