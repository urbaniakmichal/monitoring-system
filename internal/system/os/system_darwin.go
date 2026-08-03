//go:build darwin

package system

import (
	"os"
	"runtime"
)

func (*SystemOsInfo)RetrieveSystemInfo() (SystemInformation, error) {
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