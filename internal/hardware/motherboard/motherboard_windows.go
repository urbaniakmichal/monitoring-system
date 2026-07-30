//go:build windows

package motherboard

import (
	"log/slog"
	"os/exec"
	"strings"
)

func RetrieveMotherboardInfo() (MotherboardInformation, error) {
	cmd := exec.Command("powershell", "-Command", "Get-CimInstance Win32_BaseBoard | Select-Object Manufacturer, Product, Version | ConvertTo-Json")
	output, err := cmd.Output()
	if err != nil {
		slog.Error("Failed to get motherboard on Windows", slog.String("error_details", err.Error()))
		return MotherboardInformation{Manufacturer: "Unknown", Product: "Unknown"}, nil
	}

	raw := string(output)
	// Proste wyciąganie danych z JSON-a PowerShell bez zewnętrznych struktur
	return MotherboardInformation{
		Manufacturer: extractJSONField(raw, "Manufacturer"),
		Product:      extractJSONField(raw, "Product"),
		Version:      extractJSONField(raw, "Version"),
	}, nil
}

func extractJSONField(jsonStr, field string) string {
	key := `"` + field + `":`
	idx := strings.Index(jsonStr, key)
	if idx == -1 {
		return "Unknown"
	}
	start := idx + len(key)
	sub := strings.TrimSpace(jsonStr[start:])
	sub = strings.TrimPrefix(sub, `"`)
	end := strings.Index(sub, `"`)
	if end == -1 {
		return "Unknown"
	}
	return sub[:end]
}