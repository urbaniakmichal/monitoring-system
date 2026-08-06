//go:build darwin

package peripherals

import (
	"os/exec"
	"strings"
)

func (*HardwarePeripherals)RetrievePeripheralsInfo() (PeripheralsInformation, error) {
	var displays []DisplayInformation

	cmd := exec.Command("system_profiler", "SPDisplaysDataType")
	if output, err := cmd.Output(); err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "Resolution:") {
				displays = append(displays, DisplayInformation{
					Name:       "Apple Display",
					Resolution: strings.TrimPrefix(trimmed, "Resolution:"),
					IsPrimary:  len(displays) == 0,
				})
			}
		}
	}

	var usbDevices []USBDeviceInformation
	cmdUSB := exec.Command("system_profiler", "SPUSBDataType")
	if output, err := cmdUSB.Output(); err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if strings.Contains(trimmed, ":") && !strings.HasPrefix(trimmed, "USB") {
				name := strings.Split(trimmed, ":")[0]
				usbDevices = append(usbDevices, USBDeviceInformation{
					Name: name,
				})
			}
		}
	}

	return PeripheralsInformation{
		Displays: displays,
		USB:      usbDevices,
	}, nil
}