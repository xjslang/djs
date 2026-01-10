package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		os.Exit(printUsage())
	}

	switch args[0] {
	case "fmt", "format":
		fmt.Println("format")
		os.Exit(0)
	case "v", "version":
		os.Exit(printVersion())
	case "h", "help":
		os.Exit(printHelp())
	default:
		fmt.Println("compile")
		os.Exit(0)
	}
}

func printVersion() int {
	fmt.Println("version")
	return 0
}

func printHelp() int {
	fmt.Println("help")
	return 0
}

func printUsage() int {
	fmt.Println("usage")
	return 0
}
