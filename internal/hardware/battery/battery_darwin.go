//go:build darwin

package battery

import (
	"os/exec"
	"strconv"
	"strings"
)

func  (*HardwareBattery)RetrieveBatteryInfo() ([]BatteryInformation, error) {
	cmd := exec.Command("pmset", "-g", "batt")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var batteries []BatteryInformation
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if !strings.Contains(line, "InternalBattery") {
			continue
		}

		parts := strings.Split(line, "\t")
		if len(parts) < 2 {
			continue
		}

		infoPart := parts[1]
		subParts := strings.Split(infoPart, ";")
		if len(subParts) < 2 {
			continue
		}

		percentStr := strings.TrimSpace(subParts[0])
		percentStr = strings.TrimSuffix(percentStr, "%")
		percent, _ := strconv.ParseFloat(percentStr, 64)

		status := strings.TrimSpace(subParts[1])
		isCharging := strings.Contains(strings.ToLower(status), "charging") && !strings.Contains(strings.ToLower(status), "discharging")

		batteries = append(batteries, BatteryInformation{
			Name:       "InternalBattery",
			Percent:    percent,
			IsCharging: isCharging,
			Status:     status,
		})
	}

	return batteries, nil
}