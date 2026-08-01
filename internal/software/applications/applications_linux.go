//go:build linux

package applications

import (
	"bytes"
	"log/slog"
	"os/exec"
	"strings"
)

func (*SoftwareApplications) RetrieveInstalledApplications() ([]ApplicationInformation, error) {
	// Try dpkg-query (Debian/Ubuntu)
	cmd := exec.Command("dpkg-query", "-f=${Package}\t${Version}\t${Maintainer}\n", "-W")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		// Fallback to rpm (RHEL/Fedora)
		cmdRpm := exec.Command("rpm", "-qa", "--queryformat", "%{NAME}\t%{VERSION}\t%{VENDOR}\n")
		var outRpm bytes.Buffer
		cmdRpm.Stdout = &outRpm
		if errRpm := cmdRpm.Run(); errRpm != nil {
			slog.Error("Failed to query packages via dpkg or rpm", slog.String("error_details", err.Error()))
			return []ApplicationInformation{}, nil
		}
		out = outRpm
	}

	var apps []ApplicationInformation
	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, "\t")
		if len(parts) >= 2 {
			pub := ""
			if len(parts) >= 3 {
				pub = parts[2]
			}
			apps = append(apps, ApplicationInformation{
				DisplayName:    parts[0],
				DisplayVersion: parts[1],
				Publisher:      pub,
				InstallDate:    "",
			})
		}
	}
	return apps, nil
}