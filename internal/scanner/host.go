package scanner

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type HostInfo struct {
	Hostname string `json:"hostname"`
	OS       string `json:"os"`
	Kernel   string `json:"kernel"`
}

func ScanHost() (HostInfo, error) {
	var info HostInfo

	hostname, err := os.Hostname()
	if err != nil {
		return info, fmt.Errorf("reading hostname: %w", err)
	}
	info.Hostname = hostname

	osName, err := readOSRelease()
	if err != nil {
		return info, fmt.Errorf("reading OS info: %w", err)
	}
	info.OS = osName

	kernel, err := readKernelVersion()
	if err != nil {
		return info, fmt.Errorf("reading kernel version: %w", err)
	}
	info.Kernel = kernel

	return info, nil
}

func readOSRelease() (string, error) {
	file, err := os.Open("/etc/os-release")
	if err != nil {
		return "unknown", nil
	}
	defer file.Close()

	fields := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		key, value, found := strings.Cut(line, "=")
		if !found {
			continue
		}
		fields[key] = strings.Trim(value, "\"")
	}

	if pretty, ok := fields["PRETTY_NAME"]; ok {
		return pretty, nil
	}
	if name, ok := fields["NAME"]; ok {
		if ver, ok := fields["VERSION"]; ok {
			return name + " " + ver, nil
		}
		return name, nil
	}

	return "unknown", nil
}

func readKernelVersion() (string, error) {
	out, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return "", fmt.Errorf("running uname -r: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}
