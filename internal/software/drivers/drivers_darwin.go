//go:build darwin

package drivers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"os/exec"
)

type DriverInformation struct {
	DeviceName    string `json:"DeviceName"`
	DriverVersion string `json:"DriverVersion"`
	Manufacturer  string `json:"Manufacturer"`
	DriverName    string `json:"DriverName"`
}

type macExtensionItem struct {
	Name        string `json:"_name"`
	Version     string `json:"version"`
	ObtainedFrom string `json:"obtained_from"`
	KextVersion string `json:"kext_version"`
}

type macExtensionsResponse struct {
	SPExtensionsDataType []macExtensionItem `json:"SPExtensionsDataType"`
}

func RetrieveInstalledDrivers() ([]DriverInformation, error) {
	executableCommand := exec.Command("system_profiler", "SPExtensionsDataType", "-json")

	var outputBuffer bytes.Buffer
	executableCommand.Stdout = &outputBuffer

	executionError := executableCommand.Run()
	if executionError != nil {
		slog.Error("Failed to execute system_profiler for drivers/extensions",
			slog.String("error_details", executionError.Error()),
		)
		return nil, fmt.Errorf("system_profiler execution error: %w", executionError)
	}

	var profileResponse macExtensionsResponse
	decodingError := json.Unmarshal(outputBuffer.Bytes(), &profileResponse)
	if decodingError != nil {
		slog.Error("Failed to decode JSON for macOS drivers",
			slog.String("error_details", decodingError.Error()),
		)
		return nil, fmt.Errorf("json parsing error: %w", decodingError)
	}

	var driversList []DriverInformation
	for _, ext := range profileResponse.SPExtensionsDataType {
		version := ext.Version
		if version == "" {
			version = ext.KextVersion
		}
		driversList = append(driversList, DriverInformation{
			DeviceName:    ext.Name,
			DriverVersion: version,
			Manufacturer:  ext.ObtainedFrom,
			DriverName:    ext.Name,
		})
	}

	return driversList, nil
}