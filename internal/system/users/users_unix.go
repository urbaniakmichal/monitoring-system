//go:build linux || darwin

package users

import (
	"os/exec"
	"strings"
)

func RetrieveUsersInfo() ([]UserInformation, error) {
	cmd := exec.Command("who")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var users []UserInformation
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			host := ""
			if len(fields) >= 3 {
				host = strings.Trim(fields[2], "()")
			}
			users = append(users, UserInformation{
				Username: fields[0],
				Terminal: fields[1],
				Host:     host,
			})
		}
	}
	return users, nil
}