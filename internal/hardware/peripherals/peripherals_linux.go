//go:build linux

package peripherals

import (
	"os/exec"
	"strings"
)

func (*HardwarePeripherals)RetrievePeripheralsInfo() (PeripheralsInformation, error) {
	var displays []DisplayInformation

	cmd := exec.Command("xrandr")
	if output, err := cmd.Output(); err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, " connected") {
				parts := strings.Fields(line)
				name := parts[0]
				isPrimary := strings.Contains(line, "primary")
				res := "Unknown"
				for _, p := range parts {
					if strings.Contains(p, "x") && strings.Contains(p, "+") {
						res = strings.Split(p, "+")[0]
						break
					}
				}
				displays = append(displays, DisplayInformation{
					Name:       name,
					Resolution: res,
					IsPrimary:  isPrimary,
				})
			}
		}
	}

	var usbDevices []USBDeviceInformation
	cmdUSB := exec.Command("lsusb")
	if output, err := cmdUSB.Output(); err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.TrimSpace(line) == "" {
				continue
			}
			usbDevices = append(usbDevices, USBDeviceInformation{
				Name: line,
			})
		}
	}

	return PeripheralsInformation{
		Displays: displays,
		USB:      usbDevices,
	}, nil
}