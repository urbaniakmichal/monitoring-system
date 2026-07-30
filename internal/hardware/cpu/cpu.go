package cpu

type CPUInformation struct {
	ModelName    string  `json:"ModelName"`
	Cores        int     `json:"Cores"`
	UsagePercent float64 `json:"UsagePercent"`
}