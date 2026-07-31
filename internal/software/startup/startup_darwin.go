//go:build darwin

package startup

import (
	"bytes"
	"os/exec"
	"strings"
)

func RetrieveStartupCommands() ([]StartupCommandInformation, error) {
	executableCommand := exec.Command("launchctl", "print", "gui/501") // or user domain

	var outputBuffer bytes.Buffer
	executableCommand.Stdout = &outputBuffer

	// launchctl print might fail depending on user permissions, handle gracefully
	_ = executableCommand.Run()

	var startupList []StartupCommandInformation
	lines := strings.Split(outputBuffer.String(), "\n")
	for _, line := range lines {
		if strings.Contains(line, "service =") {
			parts := strings.Split(line, "=")
			if len(parts) > 1 {
				name := strings.TrimSpace(parts[1])
				startupList = append(startupList, StartupCommandInformation{
					Name:     name,
					Command:  name,
					Location: "LaunchAgents",
					User:     "current",
				})
			}
		}
	}

	return startupList, nil
}