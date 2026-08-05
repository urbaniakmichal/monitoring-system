//go:build linux

package bios

import (
	"os"
	"strings"
)

func (*HardwareBios)RetrieveBiosInformation() (BiosInformation, error) {
	readDmi := func(filename string) string {
		//#nosec G304
		data, err := os.ReadFile("/sys/class/dmi/id/" + filename)
		if err != nil {
			return "Unknown"
		}
		return strings.TrimSpace(string(data))
	}

	return BiosInformation{
		Manufacturer: readDmi("bios_vendor"),
		Version:      readDmi("bios_version"),
		ReleaseDate:  readDmi("bios_date"),
		SerialNumber: readDmi("product_serial"),
	}, nil
}