package integration

import (
	"fmt"
	"testing"

	"github.com/xjslang/djs/builder"
	"github.com/xjslang/djs/plugins"
	"github.com/xjslang/xjs/compiler"
	"github.com/xjslang/xjs/lexer"
)

func TestFormat(t *testing.T) {
	// input := `if (x%2 == 0) console.log('even')`
	input := `
	function hello() {
		let db = openDb() or {
			console.error('Failed to open db')
			return
		}
		defer closeDb(db)
	}`
	lb := lexer.NewBuilder()
	p := builder.New(lb).Build(input)
	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("ParseProgram error = %v", err)
	}
	plugins.TransformStatement(program)
	// debug.Print(program)
	result := compiler.New().WithPrettyPrint().Compile(program)
	fmt.Println(result.Code)
}
