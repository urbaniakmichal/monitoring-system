//go:build windows

package bios

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"os/exec"
)

func (*HardwareBios)RetrieveBiosInformation() (BiosInformation, error) {
	scriptContent := "Get-CimInstance Win32_BIOS | Select-Object Manufacturer, SMBIOSBIOSVersion, ReleaseDate, SerialNumber | ConvertTo-Json"

	executableCommand := exec.Command("powershell", "-NoProfile", "-Command", scriptContent)

	var outputBuffer bytes.Buffer
	executableCommand.Stdout = &outputBuffer

	if err := executableCommand.Run(); err != nil {
		slog.Error("Failed to execute PowerShell command for BIOS info", slog.String("error_details", err.Error()))
		return BiosInformation{}, fmt.Errorf("powershell execution error: %w", err)
	}

	trimmed := bytes.TrimSpace(outputBuffer.Bytes())
	if len(trimmed) == 0 {
		return BiosInformation{}, nil
	}

	type winBios struct {
		Manufacturer      string `json:"Manufacturer"`
		SMBIOSBIOSVersion string `json:"SMBIOSBIOSVersion"`
		ReleaseDate       string `json:"ReleaseDate"`
		SerialNumber      string `json:"SerialNumber"`
	}

	var wb winBios
	if err := json.Unmarshal(trimmed, &wb); err != nil {
		slog.Error("Failed to decode JSON for BIOS info", slog.String("error_details", err.Error()))
		return BiosInformation{}, fmt.Errorf("json parsing error: %w", err)
	}

	return BiosInformation{
		Manufacturer: wb.Manufacturer,
		Version:      wb.SMBIOSBIOSVersion,
		ReleaseDate:  wb.ReleaseDate,
		SerialNumber: wb.SerialNumber,
	}, nil
}