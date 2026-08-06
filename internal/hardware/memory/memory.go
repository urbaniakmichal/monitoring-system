package memory

type MemoryInformation struct {
	TotalMB     uint64  `json:"TotalMB"`
	AvailableMB uint64  `json:"AvailableMB"`
	UsedPercent float64 `json:"UsedPercent"`
}

type Memory interface {
	RetrieveMemoryInfo() (MemoryInformation, error)
}

type HardwareMemory struct{}

var _ Memory = (*HardwareMemory)(nil)