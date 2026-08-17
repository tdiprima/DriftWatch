package scanner

import (
	"fmt"
	"os/exec"
	"strings"
)

type ServiceInfo struct {
	Name string `json:"name"`
}

func ScanServices() ([]ServiceInfo, error) {
	out, err := exec.Command("systemctl", "list-units", "--type=service", "--state=running", "--no-legend", "--no-pager").Output()
	if err != nil {
		return nil, fmt.Errorf("running systemctl: %w", err)
	}

	var services []ServiceInfo
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		services = append(services, ServiceInfo{
			Name: fields[0],
		})
	}

	return services, nil
}
