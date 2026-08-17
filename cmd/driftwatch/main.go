package main

import (
	"fmt"
	"os"

	"github.com/tdiprima/driftwatch/internal/scanner"
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

func runScan() {
	fmt.Printf("DriftWatch v%s\n\n", version)

	host, err := scanner.ScanHost()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Host:   %s\n", host.Hostname)
	fmt.Printf("OS:     %s\n", host.OS)
	fmt.Printf("Kernel: %s\n", host.Kernel)

	disks, err := scanner.ScanDisks()
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: disk scan failed: %v\n", err)
	} else {
		fmt.Println("\nDisk Usage:")
		for _, disk := range disks {
			fmt.Printf("  %-20s %d%%\n", disk.MountPoint, disk.UsagePercent)
		}
	}

	ports, err := scanner.ScanPorts()
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: port scan failed: %v\n", err)
	} else {
		fmt.Println("\nListening Ports:")
		for _, port := range ports {
			fmt.Printf("  %s (%s)\n", port.Address, port.Protocol)
		}
	}

	users, err := scanner.ScanUsers()
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: user scan failed: %v\n", err)
	} else {
		fmt.Println("\nLocal Users:")
		for _, user := range users {
			fmt.Printf("  %s (uid %d)\n", user.Username, user.UID)
		}
	}

	services, err := scanner.ScanServices()
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: service scan failed: %v\n", err)
	} else {
		fmt.Println("\nRunning Services:")
		for _, svc := range services {
			fmt.Printf("  %s\n", svc.Name)
		}
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
