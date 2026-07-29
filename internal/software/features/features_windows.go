//go:build windows

package features

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"os/exec"
)

type OptionalFeatureInformation struct {
	FeatureName string `json:"FeatureName"`
	State       string `json:"State"`
}

func RetrieveOptionalFeatures() ([]OptionalFeatureInformation, error) {
	scriptContent := "Get-WindowsOptionalFeature -Online | Select-Object FeatureName, State | ConvertTo-Json"

	executableCommand := exec.Command("powershell", "-NoProfile", "-Command", scriptContent)

	var outputBuffer bytes.Buffer
	executableCommand.Stdout = &outputBuffer

	executionError := executableCommand.Run()
	if executionError != nil {
		slog.Error("Failed to execute PowerShell command for optional features",
			slog.String("error_details", executionError.Error()),
		)
		return nil, fmt.Errorf("powershell execution error: %w", executionError)
	}

	trimmedBytes := bytes.TrimSpace(outputBuffer.Bytes())
	if len(trimmedBytes) == 0 {
		return []OptionalFeatureInformation{}, nil
	}

	if trimmedBytes[0] == '{' {
		var singleOptionalFeature OptionalFeatureInformation
		decodingError := json.Unmarshal(trimmedBytes, &singleOptionalFeature)
		if decodingError != nil {
			slog.Error("Failed to decode single JSON optional feature object",
				slog.String("error_details", decodingError.Error()),
			)
			return nil, fmt.Errorf("json parsing error for single optional feature object: %w", decodingError)
		}
		return []OptionalFeatureInformation{singleOptionalFeature}, nil
	}

	var optionalFeaturesList []OptionalFeatureInformation
	decodingError := json.Unmarshal(trimmedBytes, &optionalFeaturesList)
	if decodingError != nil {
		slog.Error("Failed to decode JSON optional features list array",
			slog.String("error_details", decodingError.Error()),
		)
		return nil, fmt.Errorf("json parsing error for optional features list: %w", decodingError)
	}

	return optionalFeaturesList, nil
}
