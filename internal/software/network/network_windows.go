//go:build windows

package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
)

type winNetworkAdapter struct {
	Caption     string   `json:"Caption"`
	Description string   `json:"Description"`
	IPAddress   []string `json:"IPAddress"`
	MACAddress  string   `json:"MACAddress"`
	DHCPEnabled bool     `json:"DHCPEnabled"`
}

func RetrieveActiveNetworkAdapters() ([]NetworkAdapterInformation, error) {
	scriptContent := "Get-CimInstance Win32_NetworkAdapterConfiguration -Filter \"IPEnabled = True\" | Select-Object Caption, Description, IPAddress, MACAddress, DHCPEnabled | ConvertTo-Json"

	executableCommand := exec.Command("powershell", "-NoProfile", "-Command", scriptContent)

	var outputBuffer bytes.Buffer
	executableCommand.Stdout = &outputBuffer

	executionError := executableCommand.Run()
	if executionError != nil {
		slog.Error("Failed to execute PowerShell command for network adapters",
			slog.String("error_details", executionError.Error()),
		)
		return nil, fmt.Errorf("powershell execution error: %w", executionError)
	}

	trimmedBytes := bytes.TrimSpace(outputBuffer.Bytes())
	if len(trimmedBytes) == 0 {
		return []NetworkAdapterInformation{}, nil
	}

	var winAdapters []winNetworkAdapter
	if trimmedBytes[0] == '{' {
		var single winNetworkAdapter
		if err := json.Unmarshal(trimmedBytes, &single); err != nil {
			return nil, err
		}
		winAdapters = []winNetworkAdapter{single}
	} else {
		if err := json.Unmarshal(trimmedBytes, &winAdapters); err != nil {
			return nil, err
		}
	}

	var adaptersList []NetworkAdapterInformation
	for _, wa := range winAdapters {
		ipStr := ""
		if len(wa.IPAddress) > 0 {
			ipStr = strings.Join(wa.IPAddress, ", ")
		}
		dhcpStr := "False"
		if wa.DHCPEnabled {
			dhcpStr = "True"
		}

		adaptersList = append(adaptersList, NetworkAdapterInformation{
			Caption:     wa.Caption,
			Description: wa.Description,
			IPAddress:   ipStr,
			MACAddress:  wa.MACAddress,
			DHCPEnabled: dhcpStr,
		})
	}

	return adaptersList, nil
}