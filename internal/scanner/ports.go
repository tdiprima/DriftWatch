package scanner

import (
	"fmt"
	"os/exec"
	"strings"
)

type PortInfo struct {
	Address  string `json:"address"`
	Protocol string `json:"protocol"`
}

func ScanPorts() ([]PortInfo, error) {
	out, err := exec.Command("ss", "-tuln").Output()
	if err != nil {
		return nil, fmt.Errorf("running ss: %w", err)
	}

	var ports []PortInfo
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")

	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}

		protocol := strings.ToLower(fields[0])
		localAddr := fields[4]

		ports = append(ports, PortInfo{
			Address:  localAddr,
			Protocol: protocol,
		})
	}

	return ports, nil
}
