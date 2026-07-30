//go:build linux

package sensors

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func RetrieveSensorsInfo() ([]SensorInformation, error) {
	var sensors []SensorInformation

	matches, err := filepath.Glob("/sys/class/thermal/thermal_zone*/temp")
	if err == nil {
		for _, path := range matches {
			//nolint:gosec // hardcoded internal system thermal path
			data, err := os.ReadFile(path)
			if err != nil {
				continue
			}
			valStr := strings.TrimSpace(string(data))
			milliC, err := strconv.ParseFloat(valStr, 64)
			if err != nil {
				continue
			}
			zoneName := filepath.Base(filepath.Dir(path))
			sensors = append(sensors, SensorInformation{
				Name:  "Thermal Zone " + zoneName,
				Value: milliC / 1000.0,
				Unit:  "°C",
			})
		}
	}

	return sensors, nil
}