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
