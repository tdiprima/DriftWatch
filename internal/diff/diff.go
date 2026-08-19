package diff

import (
	"fmt"

	"github.com/tdiprima/driftwatch/internal/scanner"
	"github.com/tdiprima/driftwatch/internal/snapshot"
)

type Change struct {
	Category string
	Type     string // "added", "removed", "changed"
	Detail   string
}

type Result struct {
	Changes []Change
}

func (r Result) Count() int {
	return len(r.Changes)
}

func Compare(baseline, current snapshot.Snapshot) Result {
	var result Result

	diffPorts(baseline.Ports, current.Ports, &result)
	diffUsers(baseline.Users, current.Users, &result)
	diffServices(baseline.Services, current.Services, &result)
	diffDisks(baseline.Disks, current.Disks, &result)

	return result
}

func diffPorts(old, new []scanner.PortInfo, result *Result) {
	oldSet := make(map[string]bool)
	for _, port := range old {
		oldSet[portKey(port)] = true
	}

	newSet := make(map[string]bool)
	for _, port := range new {
		newSet[portKey(port)] = true
	}

	for _, port := range new {
		key := portKey(port)
		if !oldSet[key] {
			result.Changes = append(result.Changes, Change{
				Category: "Listening Port",
				Type:     "added",
				Detail:   key,
			})
		}
	}

	for _, port := range old {
		key := portKey(port)
		if !newSet[key] {
			result.Changes = append(result.Changes, Change{
				Category: "Listening Port",
				Type:     "removed",
				Detail:   key,
			})
		}
	}
}

func portKey(port scanner.PortInfo) string {
	return fmt.Sprintf("%s (%s)", port.Address, port.Protocol)
}

func diffUsers(old, new []scanner.UserInfo, result *Result) {
	oldSet := make(map[string]bool)
	for _, user := range old {
		oldSet[user.Username] = true
	}

	newSet := make(map[string]bool)
	for _, user := range new {
		newSet[user.Username] = true
	}

	for _, user := range new {
		if !oldSet[user.Username] {
			result.Changes = append(result.Changes, Change{
				Category: "User",
				Type:     "added",
				Detail:   user.Username,
			})
		}
	}

	for _, user := range old {
		if !newSet[user.Username] {
			result.Changes = append(result.Changes, Change{
				Category: "User",
				Type:     "removed",
				Detail:   user.Username,
			})
		}
	}
}

func diffServices(old, new []scanner.ServiceInfo, result *Result) {
	oldSet := make(map[string]bool)
	for _, svc := range old {
		oldSet[svc.Name] = true
	}

	newSet := make(map[string]bool)
	for _, svc := range new {
		newSet[svc.Name] = true
	}

	for _, svc := range new {
		if !oldSet[svc.Name] {
			result.Changes = append(result.Changes, Change{
				Category: "Service",
				Type:     "added",
				Detail:   svc.Name,
			})
		}
	}

	for _, svc := range old {
		if !newSet[svc.Name] {
			result.Changes = append(result.Changes, Change{
				Category: "Service",
				Type:     "removed",
				Detail:   svc.Name,
			})
		}
	}
}

func diffDisks(old, new []scanner.DiskInfo, result *Result) {
	oldMap := make(map[string]int)
	for _, disk := range old {
		oldMap[disk.MountPoint] = disk.UsagePercent
	}

	for _, disk := range new {
		oldPct, existed := oldMap[disk.MountPoint]
		if !existed {
			result.Changes = append(result.Changes, Change{
				Category: "Disk",
				Type:     "added",
				Detail:   fmt.Sprintf("%s at %d%%", disk.MountPoint, disk.UsagePercent),
			})
			continue
		}

		delta := disk.UsagePercent - oldPct
		if delta != 0 {
			result.Changes = append(result.Changes, Change{
				Category: "Disk",
				Type:     "changed",
				Detail:   fmt.Sprintf("%s: %d%% -> %d%%", disk.MountPoint, oldPct, disk.UsagePercent),
			})
		}
		delete(oldMap, disk.MountPoint)
	}

	for mount, pct := range oldMap {
		result.Changes = append(result.Changes, Change{
			Category: "Disk",
			Type:     "removed",
			Detail:   fmt.Sprintf("%s was %d%%", mount, pct),
		})
	}
}
