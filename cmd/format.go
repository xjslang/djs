package cmd

import (
	"flag"
	"fmt"
	"os"
)

func Format(args []string) int {
	flag.CommandLine.Parse(args)

	inputPath := flag.Arg(0)
	inputCode, err := os.ReadFile(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", inputPath)
		return 1
	}

	fmt.Println(string(inputCode))
	return 0
}
