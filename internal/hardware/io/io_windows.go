//go:build windows

package io

import (
	"encoding/json"
	"os/exec"
)

type winNetStat struct {
	InterfaceAlias string `json:"InterfaceAlias"`
	SentBytes      uint64 `json:"SentBytes"`
	ReceivedBytes  uint64 `json:"ReceivedBytes"`
}

func (*HardwareIo)RetrieveIOStats() (IOStatistics, error) {
	netStats := make(map[string]NetworkIO)
	cmd := exec.Command("powershell", "-Command", "Get-NetAdapterStatistics | Select-Object InterfaceAlias, SentBytes, ReceivedBytes | ConvertTo-Json")
	if output, err := cmd.Output(); err == nil {
		var raw []winNetStat
		trimmed := string(output)
		if len(trimmed) > 0 && trimmed[0] == '{' {
			var s winNetStat
			if json.Unmarshal([]byte(trimmed), &s) == nil {
				raw = append(raw, s)
			}
		} else {
			_ = json.Unmarshal([]byte(trimmed), &raw)
		}
		for _, s := range raw {
			netStats[s.InterfaceAlias] = NetworkIO{
				BytesSent: s.SentBytes,
				BytesRecv: s.ReceivedBytes,
			}
		}
	}

	return IOStatistics{Network: netStats, Disk: make(map[string]DiskIO)}, nil
}