package plugins

import (
	"testing"

	"github.com/xjslang/xjs/compiler"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
)

func TestOrExpressionWithLet(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "let with or and error parameter",
			input: `let db = openDatabase('mydata.db') or |err| {
				console.error('Unable to open database: ', err)
				return
			}`,
			expected: `let db;try{db=openDatabase("mydata.db")}catch(err){console.error("Unable to open database: ",err);return}`,
		},
		{
			name: "let with or without error parameter",
			input: `let result = riskyOperation() or {
				console.log('Operation failed')
				return
			}`,
			expected: `let result;try{result=riskyOperation()}catch{console.log("Operation failed");return}`,
		},
		{
			name: "let with or and simple fallback",
			input: `let x = getValue() or |e| {
				return
			}`,
			expected: `let x;try{x=getValue()}catch(e){return}`,
		},
		{
			name: "multiple statements with or",
			input: `let a = getA() or |err| { return }
			let b = getB() or { return }`,
			expected: `let a;try{a=getA()}catch(err){return};let b;try{b=getB()}catch{return}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lb := lexer.NewBuilder()
			p := parser.NewBuilder(lb).
				Install(OrPlugin).
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

func TestOrExpressionStatement(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "expression statement with or",
			input: `doSomething() or |err| {
				console.log('Failed:', err)
			}`,
			expected: `try{doSomething()}catch(err){console.log("Failed:",err)}`,
		},
		{
			name: "expression statement with or without parameter",
			input: `execute() or {
				console.log('Error occurred')
			}`,
			expected: `try{execute()}catch{console.log("Error occurred")}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lb := lexer.NewBuilder()
			p := parser.NewBuilder(lb).
				Install(OrPlugin).
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

func TestOrInFunction(t *testing.T) {
	input := `function connect() {
		let db = openDb('test.db') or |err| {
			console.log('Failed to open:', err)
			return
		}
		console.log('Connected')
	}`
	expected := `function connect(){let db;try{db=openDb("test.db")}catch(err){console.log("Failed to open:",err);return};console.log("Connected")}`

	lb := lexer.NewBuilder()
	p := parser.NewBuilder(lb).
		Install(OrPlugin).
		Build(input)
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if prog == nil {
		t.Fatal("Expected program to be parsed")
	}

	result := compiler.New().Compile(prog)
	if result.Code != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result.Code)
	}
}

func TestOrErrorCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "or without block",
			input: `let x = getValue() or`,
		},
		{
			name:  "or with incomplete error parameter",
			input: `let x = getValue() or |err`,
		},
		{
			name:  "or with missing identifier",
			input: `let x = getValue() or | |`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lb := lexer.NewBuilder()
			p := parser.NewBuilder(lb).
				Install(OrPlugin).
				Build(tt.input)
			_, err := p.ParseProgram()
			if err == nil {
				t.Error("Expected error but got none")
			}
		})
	}
}

func TestOrWithASI(t *testing.T) {
	input := `
	function test() {
		let x = getValue() or |err| {
			return
		}
		let y = getOther() or {
			return
		}
	}`

	lb := lexer.NewBuilder()
	p := parser.NewBuilder(lb).
		Install(OrPlugin).
		Build(input)
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("Expected no error with ASI, got: %v", err)
	}
	if prog == nil {
		t.Fatal("Expected program to be parsed")
	}
}

func TestNestedOrExpressions(t *testing.T) {
	input := `let db = openPrimary() or |err1| {
		let fallback = openSecondary() or |err2| {
			console.log('All failed')
			return
		}
		db = fallback
	}`

	lb := lexer.NewBuilder()
	p := parser.NewBuilder(lb).
		Install(OrPlugin).
		Build(input)
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if prog == nil {
		t.Fatal("Expected program to be parsed")
	}

	result := compiler.New().Compile(prog)
	// Just verify it compiles, the nested structure is complex
	if len(result.Code) == 0 {
		t.Error("Expected non-empty compiled code")
	}
}
