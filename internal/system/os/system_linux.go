//go:build linux

package system

import (
	"bytes"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func (*SystemOsInfo)RetrieveSystemInfo() (SystemInformation, error) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	uptimeBytes, _ := os.ReadFile("/proc/uptime")
	uptimeStr := strings.Fields(string(uptimeBytes))
	uptime := "N/A"
	if len(uptimeStr) > 0 {
		uptime = uptimeStr[0] + " seconds"
	}

	var processes []ProcessInformation
	cmd := exec.Command("ps", "-eo", "pid,comm,%cpu,%mem", "--sort=-%cpu", "h")
	if output, err := cmd.Output(); err == nil {
		lines := strings.Split(string(bytes.TrimSpace(output)), "\n")
		for i, line := range lines {
			if i >= 10 || strings.TrimSpace(line) == "" {
				continue
			}
			fields := strings.Fields(line)
			if len(fields) < 4 {
				continue
			}

			processes = append(processes, ProcessInformation{
				Name: fields[1],
			})
		}
	}

	return SystemInformation{
		Hostname:     hostname,
		OS:           runtime.GOOS,
		Architecture: runtime.GOARCH,
		Uptime:       uptime,
		TopProcesses: processes,
	}, nil
}