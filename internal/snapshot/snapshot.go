package snapshot

import (
	"fmt"
	"time"

	"github.com/tdiprima/driftwatch/internal/scanner"
)

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
