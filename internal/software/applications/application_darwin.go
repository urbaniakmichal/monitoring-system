//go:build darwin

package applications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"os/exec"
)

type ApplicationInformation struct {
	DisplayName    string `json:"DisplayName"`
	DisplayVersion string `json:"DisplayVersion"`
	Publisher      string `json:"Publisher"`
	InstallDate    string `json:"InstallDate"`
}

type macAppItem struct {
	Name           string `json:"_name"`
	Version        string `json:"version"`
	ObtainedFrom   string `json:"obtained_from"`
	LastModified   string `json:"last_modified"`
}

type macApplicationsResponse struct {
	SPApplicationsDataType []macAppItem `json:"SPApplicationsDataType"`
}

func RetrieveInstalledApplications() ([]ApplicationInformation, error) {
	executableCommand := exec.Command("system_profiler", "SPApplicationsDataType", "-json")

	var outputBuffer bytes.Buffer
	executableCommand.Stdout = &outputBuffer

	executionError := executableCommand.Run()
	if executionError != nil {
		slog.Error("Failed to execute system_profiler for applications",
			slog.String("error_details", executionError.Error()),
		)
		return nil, fmt.Errorf("system_profiler execution error: %w", executionError)
	}

	var profileResponse macApplicationsResponse
	decodingError := json.Unmarshal(outputBuffer.Bytes(), &profileResponse)
	if decodingError != nil {
		slog.Error("Failed to decode JSON for macOS applications",
			slog.String("error_details", decodingError.Error()),
		)
		return nil, fmt.Errorf("json parsing error: %w", decodingError)
	}

	var applicationsList []ApplicationInformation
	for _, app := range profileResponse.SPApplicationsDataType {
		applicationsList = append(applicationsList, ApplicationInformation{
			DisplayName:    app.Name,
			DisplayVersion: app.Version,
			Publisher:      app.ObtainedFrom,
			InstallDate:    app.LastModified,
		})
	}

	return applicationsList, nil
}