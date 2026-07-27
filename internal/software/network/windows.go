package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"os/exec"
)

type NetworkAdapterInformation struct {
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

	if trimmedBytes[0] == '{' {
		var singleNetworkAdapter NetworkAdapterInformation
		decodingError := json.Unmarshal(trimmedBytes, &singleNetworkAdapter)
		if decodingError != nil {
			slog.Error("Failed to decode single JSON network adapter object",
				slog.String("error_details", decodingError.Error()),
			)
			return nil, fmt.Errorf("json parsing error for single network adapter object: %w", decodingError)
		}
		return []NetworkAdapterInformation{singleNetworkAdapter}, nil
	}

	var networkAdapterList []NetworkAdapterInformation
	decodingError := json.Unmarshal(trimmedBytes, &networkAdapterList)
	if decodingError != nil {
		slog.Error("Failed to decode JSON network adapters list array",
			slog.String("error_details", decodingError.Error()),
		)
		return nil, fmt.Errorf("json parsing error for network adapters list: %w", decodingError)
	}

	return networkAdapterList, nil
}
