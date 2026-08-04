//go:build linux

package battery

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func  (*HardwareBattery)RetrieveBatteryInfo() ([]BatteryInformation, error) {
	var batteries []BatteryInformation
	matches, err := filepath.Glob("/sys/class/power_supply/BAT*")
	if err != nil {
		return nil, err
	}

	for _, batPath := range matches {
		name := filepath.Base(batPath)

		//nolint:gosec // internal system power supply path
		capData, err := os.ReadFile(filepath.Join(batPath, "capacity"))
		var percent float64
		if err == nil {
			percent, _ = strconv.ParseFloat(strings.TrimSpace(string(capData)), 64)
		}

		//nolint:gosec // internal system power supply path
		statusData, err := os.ReadFile(filepath.Join(batPath, "status"))
		status := "Unknown"
		isCharging := false
		if err == nil {
			status = strings.TrimSpace(string(statusData))
			isCharging = strings.EqualFold(status, "Charging")
		}

		batteries = append(batteries, BatteryInformation{
			Name:       name,
			Percent:    percent,
			IsCharging: isCharging,
			Status:     status,
		})
	}
	return batteries, nil
}