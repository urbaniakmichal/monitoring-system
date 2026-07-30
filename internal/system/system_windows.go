//go:build windows

package system

import (
	"os"
	"runtime"
)

func RetrieveSystemInfo() (SystemInformation, error) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	return SystemInformation{
		Hostname:     hostname,
		OS:           runtime.GOOS,
		Architecture: runtime.GOARCH,
		Uptime:       "Active",
		TopProcesses: []ProcessInformation{},
	}, nil
}