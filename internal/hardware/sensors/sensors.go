package sensors

type SensorInformation struct {
	Name  string  `json:"Name"`
	Value float64 `json:"Value"`
	Unit  string  `json:"Unit"`
}

type Sensors interface {
	RetrieveSensorsInfo() ([]SensorInformation, error)
}

type HardwareSensors struct{}

var _ Sensors = (*HardwareSensors)(nil)