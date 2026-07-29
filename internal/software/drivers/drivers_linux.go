//go:build linux

package drivers

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

type DriverInformation struct {
	DeviceName    string `json:"DeviceName"`
	DriverVersion string `json:"DriverVersion"`
	Manufacturer  string `json:"Manufacturer"`
	DriverName    string `json:"DriverName"`
}

func RetrieveInstalledDrivers() ([]DriverInformation, error) {
	file, err := os.Open("/proc/modules")
	if err != nil {
		slog.Error("Failed to open /proc/modules", slog.String("error_details", err.Error()))
		return nil, fmt.Errorf("failed to read modules: %w", err)
	}
	defer file.Close()

	var drivers []DriverInformation
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) > 0 {
			name := fields[0]
			drivers = append(drivers, DriverInformation{
				DeviceName:    name,
				DriverVersion: "",
				Manufacturer:  "Linux Kernel",
				DriverName:    name,
			})
		}
	}
	return drivers, nil
}