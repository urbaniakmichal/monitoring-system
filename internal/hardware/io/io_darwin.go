//go:build darwin

package io

import (
	"os/exec"
	"strconv"
	"strings"
)

func RetrieveIOStats() (IOStatistics, error) {
	netStats := make(map[string]NetworkIO)
	cmd := exec.Command("netstat", "-I", "en0", "-b")
	if output, err := cmd.Output(); err == nil {
		lines := strings.Split(string(output), "\n")
		if len(lines) > 1 {
			fields := strings.Fields(lines[1])
			if len(fields) >= 10 {
				recv, _ := strconv.ParseUint(fields[6], 10, 64)
				sent, _ := strconv.ParseUint(fields[9], 10, 64)
				netStats["en0"] = NetworkIO{BytesSent: sent, BytesRecv: recv}
			}
		}
	}
	return IOStatistics{Network: netStats, Disk: make(map[string]DiskIO)}, nil
}