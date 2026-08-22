//go:build js && wasm

package main

import (
	"syscall/js"

	"github.com/xjslang/djs"
	"github.com/xjslang/xjs/printer"
)

var (
	compileFn js.Func
	formatFn  js.Func
)

func main() {
	compileFn = js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) < 1 {
			return map[string]any{"error": "missing input"}
		}
		input := args[0].String()
		result, err := djs.Parse(input)
		if err != nil {
			return map[string]any{"error": err.Error()}
		}
		code, err := djs.Compile(result)
		if err != nil {
			return map[string]any{"error": err.Error()}
		}
		return map[string]any{"code": code}
	})
	formatFn = js.FuncOf(func(this js.Value, args []js.Value) any {
		input := args[0].String()
		result, err := djs.Parse(input)
		if err != nil {
			return map[string]any{"error": err.Error()}
		}
		code, err := djs.Format(result, printer.WithIndent("  "))
		if err != nil {
			return map[string]any{"error": err.Error()}
		}
		return map[string]any{"code": code}
	})

	js.Global().Set("djslang", js.ValueOf(map[string]any{
		"compile": compileFn,
		"format":  formatFn,
	}))

	// keep the program running
	select {}
}
