package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/xjslang/djs/djs"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/token"
)

func main() {
	var check, format bool
	flag.BoolVar(&format, "format", false, "format code")
	flag.BoolVar(&check, "check", false, "check code")
	flag.Parse()

	source := flag.Arg(0)
	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s <file.djs>\n", os.Args[0])
		os.Exit(2)
	}

	if check && format {
		fmt.Fprintf(os.Stderr, "-check and -format cannot be used together\n")
		os.Exit(2)
	}

	input, err := os.ReadFile(source)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	switch {
	case format:
		result, err := djs.Parse(input)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		code, err := djs.Format(result)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println(code)
	case check:
		_, err := djs.Parse(input)
		if errList, ok := err.(parser.ErrorList); err == nil || ok {
			fmt.Fprintf(os.Stdout, "[")
			for i, e := range errList {
				var start, end token.Position
				var msg, code string
				if pe, ok := e.(parser.Error); ok {
					start, end = pe.Range.Start, pe.Range.End
					msg = pe.Message
					code = "SYNTAX"
				} else {
					msg = e.Error()
					code = "FATAL"
				}
				if i == 0 {
					fmt.Fprint(os.Stdout, "\n")
				}
				fmt.Fprint(os.Stdout, "\t{range: {")
				fmt.Fprintf(os.Stdout, "start: {line: %d, colum: %d}, ", start.Line, start.Column)
				fmt.Fprintf(os.Stdout, "end: {line: %d, colum: %d}}}, ", end.Line, end.Column)
				fmt.Fprintf(os.Stdout, "message: %q, ", msg)
				fmt.Fprintf(os.Stdout, "code: %q}\n", code)
			}
			fmt.Fprintf(os.Stdout, "]\n")
		} else {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		result, err := djs.Parse(input)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		code, err := djs.Compile(result)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println(code)
	}
}
