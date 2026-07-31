//go:build linux

package operatingsystem

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
)

func RetrieveOperatingSystemInformation() (OperatingSystemInformation, error) {
	file, err := os.Open("/etc/os-release")
	if err != nil {
		slog.Error("Failed to open /etc/os-release", slog.String("error_details", err.Error()))
		return OperatingSystemInformation{}, fmt.Errorf("failed to open os-release: %w", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			slog.Error("Failed to close file", slog.String("error_details", err.Error()))
		}
	}()

	var name, version string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "PRETTY_NAME=") {
			name = strings.Trim(strings.TrimPrefix(line, "PRETTY_NAME="), "\"")
		} else if strings.HasPrefix(line, "VERSION_ID=") {
			version = strings.Trim(strings.TrimPrefix(line, "VERSION_ID="), "\"")
		}
	}

	archCmd := exec.Command("uname", "-m")
	archOut, _ := archCmd.Output()
	arch := strings.TrimSpace(string(archOut))

	kernelCmd := exec.Command("uname", "-r")
	kernelOut, _ := kernelCmd.Output()
	kernel := strings.TrimSpace(string(kernelOut))

	return OperatingSystemInformation{
		Caption:        name,
		Version:        version,
		BuildNumber:    kernel,
		Manufacturer:   "Linux Community",
		OSArchitecture: arch,
		InstallDate:    "",
		LastBootTime:   kernel,
	}, nil
}