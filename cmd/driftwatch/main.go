package main

import (
	"fmt"
	"os"

	"github.com/tdiprima/driftwatch/internal/diff"
	"github.com/tdiprima/driftwatch/internal/snapshot"
)

const version = "0.1.0"

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
	colorDim    = "\033[2m"
)

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

func printHeader() {
	fmt.Printf("%s%sDriftWatch%s v%s\n\n", colorBold, colorCyan, colorReset, version)
}

func runScan() {
	printHeader()

	fmt.Printf("Scanning host: ")
	snap, err := snapshot.Capture()
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n%serror: %v%s\n", colorRed, err, colorReset)
		os.Exit(1)
	}
	fmt.Printf("%s%s%s\n\n", colorBold, snap.Host.Hostname, colorReset)

	printSnapshot(snap)
}

func runDiff() {
	baseline, err := snapshot.LoadBaseline()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%serror: %v%s\n", colorRed, err, colorReset)
		fmt.Fprintln(os.Stderr, "Run 'driftwatch baseline' first to create a baseline.")
		os.Exit(1)
	}

	current, err := snapshot.Capture()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%serror: %v%s\n", colorRed, err, colorReset)
		os.Exit(1)
	}

	result := diff.Compare(baseline, current)

	if result.Count() == 0 {
		fmt.Printf("%s%sDriftWatch%s %s— No drift detected.%s\n", colorBold, colorCyan, colorReset, colorGreen, colorReset)
		fmt.Println("")
		fmt.Printf("System matches baseline. %s✓%s\n", colorGreen, colorReset)
		return
	}

	fmt.Printf("%s%sDriftWatch%s — %sChanges detected%s\n\n", colorBold, colorCyan, colorReset, colorYellow, colorReset)

	for _, change := range result.Changes {
		switch change.Type {
		case "added":
			fmt.Printf("%s+ NEW %s%s\n", colorGreen, change.Category, colorReset)
			fmt.Printf("  %s\n\n", change.Detail)
		case "removed":
			fmt.Printf("%s- %s REMOVED%s\n", colorRed, change.Category, colorReset)
			fmt.Printf("  %s\n\n", change.Detail)
		case "changed":
			fmt.Printf("%s~ %s CHANGE%s\n", colorYellow, change.Category, colorReset)
			fmt.Printf("  %s\n\n", change.Detail)
		}
	}

	fmt.Printf("%s%d change(s) detected.%s\n", colorYellow, result.Count(), colorReset)
}

func runBaseline() {
	printHeader()

	fmt.Printf("Creating baseline for %s", colorBold)
	snap, err := snapshot.Capture()
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n%serror: %v%s\n", colorRed, err, colorReset)
		os.Exit(1)
	}
	fmt.Printf("%s%s...\n\n", snap.Host.Hostname, colorReset)

	checks := []string{"Host information", "Disk usage", "Listening ports", "Local users", "Running services"}
	for _, check := range checks {
		fmt.Printf("  %s✓%s %s\n", colorGreen, colorReset, check)
	}

	path, err := snapshot.Save(snap)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n%serror: %v%s\n", colorRed, err, colorReset)
		os.Exit(1)
	}

	fmt.Printf("\n%sBaseline created successfully.%s\n", colorGreen, colorReset)
	fmt.Printf("%s%s%s\n", colorDim, path, colorReset)
}

func printSnapshot(snap snapshot.Snapshot) {
	fmt.Printf("  %sHost:%s   %s\n", colorDim, colorReset, snap.Host.Hostname)
	fmt.Printf("  %sOS:%s     %s\n", colorDim, colorReset, snap.Host.OS)
	fmt.Printf("  %sKernel:%s %s\n", colorDim, colorReset, snap.Host.Kernel)

	if len(snap.Disks) > 0 {
		fmt.Printf("\n%s%sDisk Usage%s\n", colorBold, colorCyan, colorReset)
		for _, disk := range snap.Disks {
			color := colorGreen
			if disk.UsagePercent >= 90 {
				color = colorRed
			} else if disk.UsagePercent >= 75 {
				color = colorYellow
			}
			fmt.Printf("  %-20s %s%d%%%s\n", disk.MountPoint, color, disk.UsagePercent, colorReset)
		}
	}

	if len(snap.Ports) > 0 {
		fmt.Printf("\n%s%sListening Ports%s\n", colorBold, colorCyan, colorReset)
		for _, port := range snap.Ports {
			fmt.Printf("  %s %s(%s)%s\n", port.Address, colorDim, port.Protocol, colorReset)
		}
	}

	if len(snap.Users) > 0 {
		fmt.Printf("\n%s%sLocal Users%s\n", colorBold, colorCyan, colorReset)
		for _, user := range snap.Users {
			fmt.Printf("  %-16s %s(uid %d)%s\n", user.Username, colorDim, user.UID, colorReset)
		}
	}

	if len(snap.Services) > 0 {
		fmt.Printf("\n%s%sRunning Services%s\n", colorBold, colorCyan, colorReset)
		for _, svc := range snap.Services {
			fmt.Printf("  %s\n", svc.Name)
		}
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "%s%sDriftWatch%s v%s\n\n", colorBold, colorCyan, colorReset, version)
	fmt.Fprintln(os.Stderr, "Usage: driftwatch <command>")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Commands:")
	fmt.Fprintf(os.Stderr, "  %sbaseline%s    Capture the trusted baseline snapshot\n", colorBold, colorReset)
	fmt.Fprintf(os.Stderr, "  %sscan%s        Display current system state\n", colorBold, colorReset)
	fmt.Fprintf(os.Stderr, "  %sdiff%s        Compare current state to baseline\n", colorBold, colorReset)
	fmt.Fprintf(os.Stderr, "  %sversion%s     Print version information\n", colorBold, colorReset)
}
