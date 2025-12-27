package plugins

import (
	"testing"

	"github.com/xjslang/xjs/compiler"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
)

func TestStrictEquality(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "double equals to triple equals",
			input:    `let result = a == b`,
			expected: `let result=(a===b)`,
		},
		{
			name:     "not equals to strict not equals",
			input:    `let result = a != b`,
			expected: `let result=(a!==b)`,
		},
		{
			name:     "multiple equality checks",
			input:    `let check = x == y && a != b`,
			expected: `let check=((x===y)&&(a!==b))`,
		},
		{
			name: "equality in if statement",
			input: `if (name == "test") {
				console.log("matched")
			}`,
			expected: `if ((name==="test")){console.log("matched")}`,
		},
		{
			name: "inequality in if statement",
			input: `if (status != 200) {
				console.log("error")
			}`,
			expected: `if ((status!==200)){console.log("error")}`,
		},
		{
			name:     "equality with null",
			input:    `let isNull = value == null`,
			expected: `let isNull=(value===null)`,
		},
		{
			name:     "inequality with undefined",
			input:    `let isDefined = x != undefined`,
			expected: `let isDefined=(x!==undefined)`,
		},
		{
			name: "equality in while loop",
			input: `while (count == 10) {
				count = count + 1
			}`,
			expected: `while ((count===10)){count=(count+1)}`,
		},
		{
			name: "mixed equalities",
			input: `let a = x == 1
			let b = y != 2
			let c = z == 3`,
			expected: `let a=(x===1);let b=(y!==2);let c=(z===3)`,
		},
		{
			name:     "chained comparisons",
			input:    `let check = a == b == c`,
			expected: `let check=((a===b)===c)`,
		},
		{
			name: "equality in function",
			input: `function compare(x, y) {
				return x == y
			}`,
			expected: `function compare(x,y){return (x===y)}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lb := lexer.NewBuilder()
			p := parser.NewBuilder(lb).
				Install(StrictEqualityPlugin).
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

func TestStrictEqualityWithASI(t *testing.T) {
	input := `
	function test() {
		let a = x == 1
		let b = y != 2
		return a
	}`

	lb := lexer.NewBuilder()
	p := parser.NewBuilder(lb).
		Install(StrictEqualityPlugin).
		Build(input)
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("Expected no error with ASI, got: %v", err)
	}
	if prog == nil {
		t.Fatal("Expected program to be parsed")
	}

	result := compiler.New().Compile(prog)
	if result.Code != `function test(){let a=(x===1);let b=(y!==2);return a}` {
		t.Errorf("Unexpected output: %s", result.Code)
	}
}

func TestStrictEqualityComplexExpressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "equality with string concatenation",
			input:    `let check = name == "John" + " Doe"`,
			expected: `let check=(name===("John"+" Doe"))`,
		},
		{
			name:     "inequality with arithmetic",
			input:    `let result = total != price + tax`,
			expected: `let result=(total!==(price+tax))`,
		},
		{
			name:     "equality with member access",
			input:    `let same = user.name == account.name`,
			expected: `let same=(user.name===account.name)`,
		},
		{
			name:     "inequality with function call",
			input:    `let different = getValue() != getOther()`,
			expected: `let different=(getValue()!==getOther())`,
		},
		{
			name:     "nested equalities in logical expression",
			input:    `let valid = a == b || c != d && e == f`,
			expected: `let valid=((a===b)||((c!==d)&&(e===f)))`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lb := lexer.NewBuilder()
			p := parser.NewBuilder(lb).
				Install(StrictEqualityPlugin).
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
