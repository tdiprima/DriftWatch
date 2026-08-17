package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/tdiprima/driftwatch/internal/scanner"
)

const baselineFileName = "baseline.json"

type Snapshot struct {
	Timestamp time.Time            `json:"timestamp"`
	Host      scanner.HostInfo     `json:"host"`
	Disks     []scanner.DiskInfo   `json:"disks"`
	Ports     []scanner.PortInfo   `json:"ports"`
	Users     []scanner.UserInfo   `json:"users"`
	Services  []scanner.ServiceInfo `json:"services"`
}

func Capture() (Snapshot, error) {
	snap := Snapshot{
		Timestamp: time.Now(),
	}

	host, err := scanner.ScanHost()
	if err != nil {
		return snap, fmt.Errorf("host scan: %w", err)
	}
	snap.Host = host

	disks, err := scanner.ScanDisks()
	if err != nil {
		return snap, fmt.Errorf("disk scan: %w", err)
	}
	snap.Disks = disks

	ports, err := scanner.ScanPorts()
	if err != nil {
		return snap, fmt.Errorf("port scan: %w", err)
	}
	snap.Ports = ports

	users, err := scanner.ScanUsers()
	if err != nil {
		return snap, fmt.Errorf("user scan: %w", err)
	}
	snap.Users = users

	services, err := scanner.ScanServices()
	if err != nil {
		return snap, fmt.Errorf("service scan: %w", err)
	}
	snap.Services = services

	return snap, nil
}

func dataDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("finding home directory: %w", err)
	}
	return filepath.Join(home, ".driftwatch"), nil
}

func Save(snap Snapshot) (string, error) {
	dir, err := dataDir()
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", fmt.Errorf("creating data directory: %w", err)
	}

	path := filepath.Join(dir, baselineFileName)

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshaling snapshot: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return "", fmt.Errorf("writing baseline: %w", err)
	}

	return path, nil
}

func LoadBaseline() (Snapshot, error) {
	dir, err := dataDir()
	if err != nil {
		return Snapshot{}, err
	}

	path := filepath.Join(dir, baselineFileName)

	data, err := os.ReadFile(path)
	if err != nil {
		return Snapshot{}, fmt.Errorf("reading baseline: %w", err)
	}

	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return Snapshot{}, fmt.Errorf("parsing baseline: %w", err)
	}

	return snap, nil
}
