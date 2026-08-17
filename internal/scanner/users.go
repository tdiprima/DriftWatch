package scanner

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type UserInfo struct {
	Username string `json:"username"`
	UID      int    `json:"uid"`
}

func ScanUsers() ([]UserInfo, error) {
	file, err := os.Open("/etc/passwd")
	if err != nil {
		return nil, fmt.Errorf("reading /etc/passwd: %w", err)
	}
	defer file.Close()

	var users []UserInfo
	s := bufio.NewScanner(file)

	for s.Scan() {
		line := s.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, ":")
		if len(parts) < 3 {
			continue
		}

		uid, err := strconv.Atoi(parts[2])
		if err != nil {
			continue
		}

		// Skip system accounts (UID < 1000), keep root
		if uid >= 1000 || uid == 0 {
			users = append(users, UserInfo{
				Username: parts[0],
				UID:      uid,
			})
		}
	}

	return users, s.Err()
}
