package drivers

type DriverInformation struct {
	DeviceName    string `json:"DeviceName"`
	DriverVersion string `json:"DriverVersion"`
	Manufacturer  string `json:"Manufacturer"`
	DriverName    string `json:"DriverName"`
}

type Drivers interface {
	RetrieveInstalledDrivers() ([]DriverInformation, error)
}

type SoftwareDrivers struct{}

var _ Drivers = (*SoftwareDrivers)(nil)