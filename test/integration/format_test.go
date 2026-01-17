package integration

import (
	"testing"

	"github.com/xjslang/djs/builder"
	"github.com/xjslang/djs/formatters"
	"github.com/xjslang/xjs/compiler"
	"github.com/xjslang/xjs/lexer"
)

func TestFormat(t *testing.T) {
	input := `function processData(items){let result=[];for(let i=0;i<items.length;i++){let item=items[i];if(item.active){result.push(item.value);}}return result;}async function main(){let db=connect("localhost") or |err|{console.error("Connection failed:",err);return;};defer db.close();let users=await db.query("SELECT * FROM users") or {console.log("Query failed");return [];};for(let i=0;i<users.length;i++){let user=users[i];if(user.age>=18){console.log(user.name);}}defer console.log("Cleanup complete");}`

	lb := lexer.NewBuilder()
	p := builder.New(lb).Build(input)
	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("ParseProgram error = %v", err)
	}

	formatters.Prepare(program)
	result := compiler.New().WithPrettyPrint(compiler.WithSemi(false)).Compile(program)

	expected := `function processData(items) {
  let result = []
  for (let i = 0; i < items.length; i++) {
    let item = items[i]
    if (item.active) {
      result.push(item.value)
    }
  }
  return result
}
async function main() {
  let db = connect("localhost") or |err| {
    console.error("Connection failed:", err)
    return
  }
  defer db.close()
  let users = await db.query("SELECT * FROM users") or {
    console.log("Query failed")
    return []
  }
  for (let i = 0; i < users.length; i++) {
    let user = users[i]
    if (user.age >= 18) {
      console.log(user.name)
    }
  }
  defer console.log("Cleanup complete")
}`

	if result.Code != expected {
		t.Errorf("Format output mismatch\nGot:\n%s\n\nExpected:\n%s", result.Code, expected)
	}
}
