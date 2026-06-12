package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/xjslang/djs/djs"
)

func main() {
	var format bool
	flag.BoolVar(&format, "format", false, "format code")
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

	var code string
	if format {
		code, err = djs.Format(result)
	} else {
		code, err = djs.Compile(result)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(code)
}
