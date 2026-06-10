package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/xjslang/djs/djs"
)

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s <file.djs>\n", os.Args[0])
		os.Exit(2)
	}

	source := flag.Arg(0)
	input, err := os.ReadFile(source)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// transform input into AST (Abstract Syntax Tree)
	result, err := djs.Parse(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// transform AST into a compiled string
	compiled := djs.Compile(result)
	fmt.Println(compiled)
}
