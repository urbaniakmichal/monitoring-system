//go:build darwin

package bios

import (
	"bytes"
	"os/exec"
	"strings"
)

func (*HardwareBios)RetrieveBiosInformation() (BiosInformation, error) {
	cmd := exec.Command("system_profiler", "SPHardwareDataType")
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return BiosInformation{}, err
	}

	info := BiosInformation{
		Manufacturer: "Apple Inc.",
		Version:      "Unknown",
		ReleaseDate:  "Unknown",
		SerialNumber: "Unknown",
	}

	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Boot ROM Version:") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				info.Version = strings.TrimSpace(parts[1])
			}
		} else if strings.HasPrefix(line, "Serial Number") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				info.SerialNumber = strings.TrimSpace(parts[1])
			}
		}
	}

	return info, nil
}