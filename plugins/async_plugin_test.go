package plugins

import (
	"testing"

	"github.com/xjslang/xjs/compiler"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
)

func TestAwaitOutsideFunction(t *testing.T) {
	input := `let x = await get()`
	lb := lexer.NewBuilder()
	p := parser.NewBuilder(lb).Install(AsyncPlugin).Build(input)
	_, err := p.ParseProgram()
	if err == nil {
		t.Errorf("Expected error when await is used outside function, but got none")
	}
}

func TestAsyncFunctionDeclaration(t *testing.T) {
	input := `async function read() { let u = await get() }`
	lb := lexer.NewBuilder()
	p := parser.NewBuilder(lb).Install(AsyncPlugin).Build(input)
	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("Parser error: %v", err)
	}
	result := compiler.New().Compile(program)
	if result.Code != "async function read(){let u=await get()}" {
		t.Errorf("Unexpected transpiled output: %q", result.Code)
	}
}

func TestAsyncFunctionExpression(t *testing.T) {
	input := `let f = async function() { await go() }`
	lb := lexer.NewBuilder()
	p := parser.NewBuilder(lb).Install(AsyncPlugin).Build(input)
	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("Parser error: %v", err)
	}
	result := compiler.New().Compile(program)
	if result.Code != "let f=async function(){await go()}" {
		t.Errorf("Unexpected transpiled output: %q", result.Code)
	}
}
