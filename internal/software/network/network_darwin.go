//go:build darwin

package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"os/exec"
)

type macNetworkItem struct {
	Name        string `json:"_name"`
	Hardware    string `json:"hardware"`
	IPv4Address string `json:"ipv4_address"`
	DHCP        string `json:"dhcp_configured"`
}

type macNetworkResponse struct {
	SPNetworkDataType []macNetworkItem `json:"SPNetworkDataType"`
}

func (*SoftwareNetworks)RetrieveActiveNetworkAdapters() ([]NetworkAdapterInformation, error) {
	executableCommand := exec.Command("system_profiler", "SPNetworkDataType", "-json")

	var outputBuffer bytes.Buffer
	executableCommand.Stdout = &outputBuffer

	executionError := executableCommand.Run()
	if executionError != nil {
		slog.Error("Failed to execute system_profiler for network adapters",
			slog.String("error_details", executionError.Error()),
		)
		return nil, fmt.Errorf("system_profiler execution error: %w", executionError)
	}

	var profileResponse macNetworkResponse
	decodingError := json.Unmarshal(outputBuffer.Bytes(), &profileResponse)
	if decodingError != nil {
		slog.Error("Failed to decode JSON for macOS network",
			slog.String("error_details", decodingError.Error()),
		)
		return nil, fmt.Errorf("json parsing error: %w", decodingError)
	}

	var adaptersList []NetworkAdapterInformation
	for _, net := range profileResponse.SPNetworkDataType {
		adaptersList = append(adaptersList, NetworkAdapterInformation{
			Description: net.Name,
			MACAddress:  net.Hardware,
			IPAddress:   net.IPv4Address,
			DHCPEnabled: net.DHCP,
		})
	}

	return adaptersList, nil
}