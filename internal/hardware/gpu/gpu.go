package gpu

type GPUInformation struct {
	Name   string `json:"Name"`
	Vendor string `json:"Vendor"`
}

type Gpu interface {
	RetrieveGPUInfo() ([]GPUInformation, error)
}

type HardwareGpu struct{}

var _ Gpu = (*HardwareGpu)(nil)