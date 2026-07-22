package metrics

import "monitor-agent/internal/system"

type Collectors struct {
	CPU       system.CPUProvider
	Memory    system.MemoryProvider
	Disk      system.DiskProvider
	Hostname  system.HostnameProvider
	Uptime    system.UptimeProvider
	Processes system.ProcessesProvider
}

func NewCollectors() Collectors {
	return Collectors{
		CPU:       system.GopsutilCPU{},
		Memory:    system.GopsutilMemory{},
		Disk:      system.GopsutilDisk{},
		Hostname:  system.GopsutilHostname{},
		Uptime:    system.GopsutilUptime{},
		Processes: system.GopsutilProcesses{},
	}
}
