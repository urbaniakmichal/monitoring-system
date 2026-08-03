package system

type ProcessInformation struct {
	PID        int32   `json:"PID"`
	Name       string  `json:"Name"`
	CPUPercent float64 `json:"CPUPercent"`
	MemoryMB   uint64  `json:"MemoryMB"`
}

type SystemInformation struct {
	Hostname     string               `json:"Hostname"`
	OS           string               `json:"OS"`
	Architecture string               `json:"Architecture"`
	Uptime       string               `json:"Uptime"`
	TopProcesses []ProcessInformation `json:"TopProcesses"`
}

type System interface {
	RetrieveSystemInfo() (SystemInformation, error)
}

type SystemOsInfo struct{}

var _ System = (*SystemOsInfo)(nil)