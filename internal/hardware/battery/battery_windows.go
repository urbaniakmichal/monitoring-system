//go:build windows

package battery

import (
	"encoding/json"
	"os/exec"
)

type winBattery struct {
	Name                     string `json:"Name"`
	EstimatedChargeRemaining uint   `json:"EstimatedChargeRemaining"`
	BatteryStatus            uint   `json:"BatteryStatus"`
}

func  (*HardwareBattery)RetrieveBatteryInfo() ([]BatteryInformation, error) {
	cmd := exec.Command("powershell", "-Command", "Get-CimInstance Win32_Battery | Select-Object Name, EstimatedChargeRemaining, BatteryStatus | ConvertTo-Json")
	output, err := cmd.Output()
	if err != nil {
		return nil, nil //Without battery
	}

	var raw []winBattery
	trimmed := string(output)
	if len(trimmed) > 0 && trimmed[0] == '{' {
		var b winBattery
		if json.Unmarshal([]byte(trimmed), &b) == nil {
			raw = append(raw, b)
		}
	} else {
		_ = json.Unmarshal([]byte(trimmed), &raw)
	}

	var batteries []BatteryInformation
	for _, b := range raw {
		isCharging := b.BatteryStatus == 2
		batteries = append(batteries, BatteryInformation{
			Name:       b.Name,
			Percent:    float64(b.EstimatedChargeRemaining),
			IsCharging: isCharging,
			Status:     "Active",
		})
	}
	return batteries, nil
}