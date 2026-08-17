package main

import (
	"fmt"
	"os"
)

const version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "version":
		fmt.Printf("DriftWatch v%s\n", version)
	case "scan":
		fmt.Println("scan: not yet implemented")
		os.Exit(1)
	case "diff":
		fmt.Println("diff: not yet implemented")
		os.Exit(1)
	case "baseline":
		fmt.Println("baseline: not yet implemented")
		os.Exit(1)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "Usage: driftwatch <command>")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Commands:")
	fmt.Fprintln(os.Stderr, "  baseline    Capture the trusted baseline snapshot")
	fmt.Fprintln(os.Stderr, "  scan        Display current system state")
	fmt.Fprintln(os.Stderr, "  diff        Compare current state to baseline")
	fmt.Fprintln(os.Stderr, "  version     Print version information")
}
