package scanner

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type DiskInfo struct {
	MountPoint   string `json:"mount_point"`
	UsagePercent int    `json:"usage_percent"`
}

func ScanDisks() ([]DiskInfo, error) {
	out, err := exec.Command("df", "--output=pcent,target", "-x", "tmpfs", "-x", "devtmpfs", "-x", "squashfs").Output()
	if err != nil {
		return nil, fmt.Errorf("running df: %w", err)
	}

	var disks []DiskInfo
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")

	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		pct, err := strconv.Atoi(strings.TrimSuffix(fields[0], "%"))
		if err != nil {
			continue
		}

		disks = append(disks, DiskInfo{
			MountPoint:   fields[1],
			UsagePercent: pct,
		})
	}

	return disks, nil
}
