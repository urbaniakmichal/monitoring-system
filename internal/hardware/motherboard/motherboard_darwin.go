//go:build darwin

package motherboard

import (
	"os/exec"
	"strings"
)

func (*HardwareMotherboard)RetrieveMotherboardInfo() (MotherboardInformation, error) {
	cmd := exec.Command("sysctl", "-n", "hw.model")
	output, err := cmd.Output()
	if err != nil {
		return MotherboardInformation{Manufacturer: "Apple", Product: "Mac"}, nil
	}

	return MotherboardInformation{
		Manufacturer: "Apple",
		Product:      strings.TrimSpace(string(output)),
		Version:      "Standard",
	}, nil
}