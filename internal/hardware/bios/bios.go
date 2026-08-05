package bios

type BiosInformation struct {
	Manufacturer string `json:"Manufacturer"`
	Version      string `json:"Version"`
	ReleaseDate  string `json:"ReleaseDate"`
	SerialNumber string `json:"SerialNumber"`
}

type Bios interface {
	RetrieveBiosInformation() (BiosInformation, error)
}

type HardwareBios struct{}

var _ Bios = (*HardwareBios)(nil)