//go:build windows

package sensors

import (
	"os/exec"
	"strconv"
	"strings"
)

func RetrieveSensorsInfo() ([]SensorInformation, error) {
	cmd := exec.Command("powershell", "-Command", "Get-CimInstance -Namespace root\\wmi -ClassName MSAcpi_ThermalZoneTemperature | Select-Object -ExpandProperty CurrentTemperature")
	output, err := cmd.Output()
	if err != nil {
		return nil, nil // WMI thermal zones might not be exposed on virtual machines
	}

	var sensors []SensorInformation
	lines := strings.Split(string(output), "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			kelvinDecide, err := strconv.ParseFloat(trimmed, 64)
			if err == nil {
				// Kelwins
				celsius := (kelvinDecide / 10.0) - 273.15
				sensors = append(sensors, SensorInformation{
					Name:  "Thermal Zone " + strconv.Itoa(i),
					Value: celsius,
					Unit:  "°C",
				})
			}
		}
	}
	return sensors, nil
}