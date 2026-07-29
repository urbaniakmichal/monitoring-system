//go:build darwin

package operatingSystem

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
)

type OperatingSystemInformation struct {
	Caption        string `json:"Caption"`
	Version        string `json:"Version"`
	BuildNumber    string `json:"BuildNumber"`
	Manufacturer   string `json:"Manufacturer"`
	OSArchitecture string `json:"OSArchitecture"`
	InstallDate    string `json:"InstallDate"`
	LastBootTime   string `json:"LastBootUpTime"`
}

type macSystemProfilerResponse struct {
	SPSoftwareDataType []struct {
		Name          string `json:"_name"`
		OSVersion     string `json:"os_version"`
		BuildVersion  string `json:"build"`
		KernelVersion string `json:"kernel_version"`
	} `json:"SPSoftwareDataType"`
}

func RetrieveOperatingSystemInformation() (OperatingSystemInformation, error) {
	executableCommand := exec.Command("system_profiler", "SPSoftwareDataType", "-json")

	var outputBuffer bytes.Buffer
	executableCommand.Stdout = &outputBuffer

	executionError := executableCommand.Run()
	if executionError != nil {
		slog.Error("Failed to execute system_profiler for OS info",
			slog.String("error_details", executionError.Error()),
		)
		return OperatingSystemInformation{}, fmt.Errorf("system_profiler execution error: %w", executionError)
	}

	var profileResponse macSystemProfilerResponse
	decodingError := json.Unmarshal(outputBuffer.Bytes(), &profileResponse)
	if decodingError != nil {
		slog.Error("Failed to decode JSON from system_profiler",
			slog.String("error_details", decodingError.Error()),
		)
		return OperatingSystemInformation{}, fmt.Errorf("json parsing error: %w", decodingError)
	}

	if len(profileResponse.SPSoftwareDataType) == 0 {
		return OperatingSystemInformation{}, fmt.Errorf("no software data returned from system_profiler")
	}

	softwareData := profileResponse.SPSoftwareDataType[0]

	archCmd := exec.Command("uname", "-m")
	archOutput, _ := archCmd.Output()
	architecture := strings.TrimSpace(string(archOutput))

	return OperatingSystemInformation{
		Caption:        fmt.Sprintf("macOS %s", softwareData.OSVersion),
		Version:        softwareData.OSVersion,
		BuildNumber:    softwareData.BuildVersion,
		Manufacturer:   "Apple Inc.",
		OSArchitecture: architecture,
		InstallDate:    "",
		LastBootTime:   softwareData.KernelVersion,
	}, nil
}