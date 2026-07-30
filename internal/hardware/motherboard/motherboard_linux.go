//go:build linux

package motherboard

import (
	"os"
	"strings"
)

func readDMIFile(path string) string {
	//nolint:gosec // path is a hardcoded internal system path for DMI info
	data, err := os.ReadFile(path)
	if err != nil {
		return "Unknown"
	}
	return strings.TrimSpace(string(data))
}

func RetrieveMotherboardInfo() (MotherboardInformation, error) {
	return MotherboardInformation{
		Manufacturer: readDMIFile("/sys/class/dmi/id/board_vendor"),
		Product:      readDMIFile("/sys/class/dmi/id/board_name"),
		Version:      readDMIFile("/sys/class/dmi/id/board_version"),
	}, nil
}