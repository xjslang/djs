package main

import (
	"fmt"
	"os"

	"github.com/xjslang/djs/cmd"
)

// Build-time injected variables
var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		os.Exit(printUsage())
	}

	switch args[0] {
	case "fmt", "format":
		fmtArgs := args[1:] // exclude the `fmt` argument
		os.Exit(cmd.Format(fmtArgs))
	case "v", "version":
		os.Exit(printVersion())
	case "h", "help":
		os.Exit(printHelp())
	default:
		os.Exit(cmd.Compile(args))
	}
}

func printVersion() int {
	fmt.Printf("djs version %s\n", Version)
	fmt.Printf("  commit: %s\n", GitCommit)
	fmt.Printf("  built:  %s\n", BuildDate)
	return 0
}

// TODO: implement help
func printHelp() int {
	fmt.Println("help")
	return 0
}

// TODO: implement usage
func printUsage() int {
	fmt.Println("usage")
	return 0
}
