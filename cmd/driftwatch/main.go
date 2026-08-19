package main

import (
	"fmt"
	"os"

	"github.com/tdiprima/driftwatch/internal/diff"
	"github.com/tdiprima/driftwatch/internal/snapshot"
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
		runScan()
	case "diff":
		runDiff()
	case "baseline":
		runBaseline()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func runScan() {
	fmt.Printf("DriftWatch v%s\n\n", version)

	snap, err := snapshot.Capture()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	printSnapshot(snap)
}

func runDiff() {
	baseline, err := snapshot.LoadBaseline()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		fmt.Fprintln(os.Stderr, "Run 'driftwatch baseline' first to create a baseline.")
		os.Exit(1)
	}

	current, err := snapshot.Capture()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	result := diff.Compare(baseline, current)

	if result.Count() == 0 {
		fmt.Println("DriftWatch — No drift detected.")
		fmt.Println("")
		fmt.Println("System matches baseline.")
		return
	}

	fmt.Println("DriftWatch — Changes detected")
	fmt.Println("")

	for _, change := range result.Changes {
		switch change.Type {
		case "added":
			fmt.Printf("+ NEW %s\n", change.Category)
			fmt.Printf("  %s\n\n", change.Detail)
		case "removed":
			fmt.Printf("- %s REMOVED\n", change.Category)
			fmt.Printf("  %s\n\n", change.Detail)
		case "changed":
			fmt.Printf("~ %s CHANGE\n", change.Category)
			fmt.Printf("  %s\n\n", change.Detail)
		}
	}

	fmt.Printf("%d change(s) detected.\n", result.Count())
}

func runBaseline() {
	fmt.Printf("DriftWatch v%s\n\n", version)
	fmt.Println("Creating baseline...")

	snap, err := snapshot.Capture()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	path, err := snapshot.Save(snap)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nBaseline saved: %s\n", path)
}

func printSnapshot(snap snapshot.Snapshot) {
	fmt.Printf("Host:   %s\n", snap.Host.Hostname)
	fmt.Printf("OS:     %s\n", snap.Host.OS)
	fmt.Printf("Kernel: %s\n", snap.Host.Kernel)

	fmt.Println("\nDisk Usage:")
	for _, disk := range snap.Disks {
		fmt.Printf("  %-20s %d%%\n", disk.MountPoint, disk.UsagePercent)
	}

	fmt.Println("\nListening Ports:")
	for _, port := range snap.Ports {
		fmt.Printf("  %s (%s)\n", port.Address, port.Protocol)
	}

	fmt.Println("\nLocal Users:")
	for _, user := range snap.Users {
		fmt.Printf("  %s (uid %d)\n", user.Username, user.UID)
	}

	fmt.Println("\nRunning Services:")
	for _, svc := range snap.Services {
		fmt.Printf("  %s\n", svc.Name)
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
