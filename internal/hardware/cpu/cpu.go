package cpu

type CPUInformation struct {
	ModelName    string  `json:"ModelName"`
	Cores        int     `json:"Cores"`
	UsagePercent float64 `json:"UsagePercent"`
}

type Cpu interface {
	RetrieveCPUInfo() (CPUInformation, error)
}

type HardwareCpu struct{}

var _ Cpu = (*HardwareCpu)(nil)