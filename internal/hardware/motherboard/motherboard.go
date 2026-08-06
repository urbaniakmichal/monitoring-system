package motherboard

type MotherboardInformation struct {
	Manufacturer string `json:"Manufacturer"`
	Product      string `json:"Product"`
	Version      string `json:"Version"`
}

type Motherboard interface {
	RetrieveMotherboardInfo() (MotherboardInformation, error)
}

type HardwareMotherboard struct{}

var _ Motherboard = (*HardwareMotherboard)(nil)