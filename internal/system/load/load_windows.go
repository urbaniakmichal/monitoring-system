//go:build windows

package load

import (
	"encoding/json"
	"os/exec"
)

type winCPU struct {
	LoadPercentage *uint32 `json:"LoadPercentage"`
}

func RetrieveLoadInfo() (LoadInformation, error) {
	cmd := exec.Command("powershell", "-Command", "Get-CimInstance Win32_Processor | Select-Object LoadPercentage | ConvertTo-Json")
	output, err := cmd.Output()
	if err != nil {
		return LoadInformation{}, nil
	}
	var cpu winCPU
	if json.Unmarshal(output, &cpu) == nil && cpu.LoadPercentage != nil {
		val := float64(*cpu.LoadPercentage)
		return LoadInformation{Load1: val, Load5: val, Load15: val}, nil
	}
	return LoadInformation{}, nil
}