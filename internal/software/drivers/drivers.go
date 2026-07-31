package drivers

type DriverInformation struct {
	DeviceName    string `json:"DeviceName"`
	DriverVersion string `json:"DriverVersion"`
	Manufacturer  string `json:"Manufacturer"`
	DriverName    string `json:"DriverName"`
}