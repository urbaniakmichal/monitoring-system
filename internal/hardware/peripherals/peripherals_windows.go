//go:build windows

package peripherals

import (
	"encoding/json"
	"os/exec"
)

type winMonitor struct {
	DeviceID string `json:"DeviceID"`
	Caption  string `json:"Caption"`
}

type winUSB struct {
	Name         string `json:"Name"`
	Manufacturer string `json:"Manufacturer"`
	DeviceID     string `json:"DeviceID"`
}

func (*HardwarePeripherals)RetrievePeripheralsInfo() (PeripheralsInformation, error) {
	var displays []DisplayInformation
	cmdMon := exec.Command("powershell", "-Command", "Get-CimInstance Win32_DesktopMonitor | Select-Object DeviceID, Caption | ConvertTo-Json")
	if output, err := cmdMon.Output(); err == nil {
		var rawMon []winMonitor
		trimmed := string(output)
		if len(trimmed) > 0 && trimmed[0] == '{' {
			var m winMonitor
			if json.Unmarshal([]byte(trimmed), &m) == nil {
				rawMon = append(rawMon, m)
			}
		} else {
			_ = json.Unmarshal([]byte(trimmed), &rawMon)
		}

		for i, m := range rawMon {
			displays = append(displays, DisplayInformation{
				Name:       m.Caption,
				Resolution: "N/A",
				IsPrimary:  i == 0,
			})
		}
	}

	var usbDevices []USBDeviceInformation
	cmdUSB := exec.Command("powershell", "-Command", "Get-CimInstance Win32_USBHub | Select-Object Name, Manufacturer, DeviceID | ConvertTo-Json")
	if output, err := cmdUSB.Output(); err == nil {
		var rawUSB []winUSB
		trimmed := string(output)
		if len(trimmed) > 0 && trimmed[0] == '{' {
			var u winUSB
			if json.Unmarshal([]byte(trimmed), &u) == nil {
				rawUSB = append(rawUSB, u)
			}
		} else {
			_ = json.Unmarshal([]byte(trimmed), &rawUSB)
		}

		for _, u := range rawUSB {
			usbDevices = append(usbDevices, USBDeviceInformation{
				Name:         u.Name,
				Manufacturer: u.Manufacturer,
				DeviceID:     u.DeviceID,
			})
		}
	}

	return PeripheralsInformation{
		Displays: displays,
		USB:      usbDevices,
	}, nil
}