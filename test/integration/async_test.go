package integration

import (
	"fmt"
	"testing"

	"github.com/xjslang/djs/builder"
	"github.com/xjslang/xjs/compiler"
	"github.com/xjslang/xjs/lexer"
)

func TestAsyncPlugin(t *testing.T) {
	input := `let x = async function main() {
		console.log('Hello, ')
		defer console.log('World!')
	}
	x()`
	lb := lexer.NewBuilder()
	p := builder.New(lb).Build(input)
	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("ParseProgram error: %v", err)
	}
	result := compiler.New().Compile(program)
	fmt.Println(result.Code)
}

func TestAsyncFunctionDeclaration(t *testing.T) {
	input := `async function read() { let u = await get() }`
	lb := lexer.NewBuilder()
	p := builder.New(lb).Build(input)
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
	p := builder.New(lb).Build(input)
	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("Parser error: %v", err)
	}
	result := compiler.New().Compile(program)
	if result.Code != "let f=async function(){await go()}" {
		t.Errorf("Unexpected transpiled output: %q", result.Code)
	}
}
