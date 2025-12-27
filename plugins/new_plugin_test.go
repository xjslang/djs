package plugins

import (
	"testing"

	"github.com/xjslang/xjs/compiler"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
)

func TestNewExpression(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "new without parentheses",
			input:    `let d = new Date`,
			expected: `let d=new Date`,
		},
		{
			name:     "new with empty parentheses",
			input:    `let d = new Date()`,
			expected: `let d=new Date()`,
		},
		{
			name:     "new with single argument",
			input:    `let e = new Error("message")`,
			expected: `let e=new Error("message")`,
		},
		{
			name:     "new with multiple arguments",
			input:    `let p = new Point(10, 20)`,
			expected: `let p=new Point(10,20)`,
		},
		{
			name:     "new with expression argument",
			input:    `let e = new Error("Error: " + msg)`,
			expected: `let e=new Error(("Error: "+msg))`,
		},
		{
			name: "new in function",
			input: `function createError() {
				return new Error("failed")
			}`,
			expected: `function createError(){return new Error("failed")}`,
		},
		{
			name:     "new with member expression",
			input:    `let c = new MyModule.MyClass()`,
			expected: `let c=new MyModule.MyClass()`,
		},
		{
			name:     "nested new expressions",
			input:    `let x = new Wrapper(new Inner())`,
			expected: `let x=new Wrapper(new Inner())`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lb := lexer.NewBuilder()
			p := parser.NewBuilder(lb).
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

func TestNewExpressionWithASI(t *testing.T) {
	// new with newline (ASI)
	input := `
	function test() {
		let e = new Error("error")
		let d = new Date
		return e
	}`
	lb := lexer.NewBuilder()
	p := parser.NewBuilder(lb).
		Install(NewPlugin).
		Build(input)
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("Expected no error with ASI, got: %v", err)
	}
	if prog == nil {
		t.Fatal("Expected program to be parsed")
	}
}
