package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"unsafe"

	"github.com/xjslang/djs"
	"github.com/xjslang/xjs/token"
)

func main() {
	var check, format, stdin bool
	flag.BoolVar(&format, "format", false, "format code")
	flag.BoolVar(&check, "check", false, "check code")
	flag.BoolVar(&stdin, "stdin", false, "read from stdin")
	flag.Parse()

	if check && format {
		fmt.Fprintf(os.Stderr, "-check and -format cannot be used together\n")
		os.Exit(2)
	}

	// read input from stdin or file
	var input []byte
	var err error
	if n := flag.NArg(); stdin && n == 0 {
		input, err = io.ReadAll(os.Stdin)
	} else if n == 1 {
		source := flag.Arg(0)
		input, err = os.ReadFile(source)
	} else {
		cmd := os.Args[0]
		fmt.Fprint(os.Stderr, "Usage:\n\n")
		fmt.Fprintf(os.Stderr, "\t%s <file.djs>\n", cmd)
		fmt.Fprintf(os.Stderr, "\t%s -format <file.djs>\n", cmd)
		fmt.Fprintf(os.Stderr, "\t%s -check <file.djs>\n", cmd)
		fmt.Fprintf(os.Stderr, "\t%s -check -stdin <<< \"DJS code\"\n", cmd)
		fmt.Fprint(os.Stderr, "\n")
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	src := unsafe.String(unsafe.SliceData(input), len(input))
	switch {
	case format:
		result, err := djs.Parse(src)
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
		_, err := djs.Parse(src)
		var errList token.ErrorList
		if err == nil || errors.As(err, &errList) {
			fmt.Fprintf(os.Stdout, "{\"errors\": [\n")
			for i, e := range errList {
				var line0, col0 int
				var line1, col1 int
				var msg, code string
				if pe, ok := errors.AsType[token.Error](e); ok {
					line0, col0 = token.Position(src, pe.Range[0])
					line1, col1 = token.Position(src, pe.Range[1])
					msg = pe.Message
					code = "SYNTAX"
				} else {
					msg = e.Error()
					code = "FATAL"
				}
				if i > 0 {
					fmt.Fprint(os.Stdout, ",\n")
				}
				fmt.Fprint(os.Stdout, "\t{\"range\": {")
				fmt.Fprintf(os.Stdout, "\"start\": {\"line\": %d, \"column\": %d}, ", line0, col0)
				fmt.Fprintf(os.Stdout, "\"end\": {\"line\": %d, \"column\": %d}}, ", line1, col1)
				fmt.Fprintf(os.Stdout, "\"message\": %q, ", msg)
				fmt.Fprintf(os.Stdout, "\"code\": %q}", code)
			}
			fmt.Fprintf(os.Stdout, "]}\n")
		} else {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		result, err := djs.Parse(src)
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
