//go:build darwin

package sensors

func (*HardwareSensors)RetrieveSensorsInfo() ([]SensorInformation, error) {
	// macOS requires external CLI tools to read temperature sensors
	return []SensorInformation{}, nil
}