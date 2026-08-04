package battery

type BatteryInformation struct {
	Name       string  `json:"Name"`
	Percent    float64 `json:"Percent"`
	IsCharging bool    `json:"IsCharging"`
	Status     string  `json:"Status"`
}

type Battery interface {
	RetrieveBatteryInfo() ([]BatteryInformation, error)
}

type HardwareBattery struct{}

var _ Battery = (*HardwareBattery)(nil)