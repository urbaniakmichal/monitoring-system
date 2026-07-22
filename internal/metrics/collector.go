package metrics

import (
	"context"
	"log/slog"
	"monitor-agent/internal/system"
	"sync"
	"time"
)

func Collect(providers Collectors, timeout time.Duration) Metrics {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var wg sync.WaitGroup

	var (
		hostname    string
		hostnameErr error
	)
	wg.Add(1)
	go func() {
		defer wg.Done()
		hostname, hostnameErr = system.Hostname(ctx, providers.Hostname)
	}()

	var (
		cpu    int
		cpuErr error
	)
	wg.Add(1)
	go func() {
		defer wg.Done()
		cpu, cpuErr = system.CPU(ctx, providers.CPU)
	}()

	var (
		memory    int
		memoryErr error
	)
	wg.Add(1)
	go func() {
		defer wg.Done()
		memory, memoryErr = system.Memory(ctx, providers.Memory)
	}()

	var (
		disk    int
		diskErr error
	)
	wg.Add(1)
	go func() {
		defer wg.Done()
		disk, diskErr = system.Disk(ctx, providers.Disk)
	}()

	var (
		uptime    string
		uptimeErr error
	)
	wg.Add(1)
	go func() {
		defer wg.Done()
		uptime, uptimeErr = system.Uptime(ctx, providers.Uptime)
	}()

	var (
		processes    int
		processesErr error
	)
	wg.Add(1)
	go func() {
		defer wg.Done()
		processes, processesErr = system.Processes(ctx, providers.Processes)
	}()

	wg.Wait()

	if hostnameErr != nil {
		slog.Warn("Failed to collect metric", "metric", "hostname", "error", hostnameErr)
		hostname = ""
	}

	if cpuErr != nil {
		slog.Warn("Failed to collect metric", "metric", "cpu", "error", cpuErr)
		cpu = 0
	}

	if memoryErr != nil {
		slog.Warn("Failed to collect metric", "metric", "memory", "error", memoryErr)
		memory = 0
	}

	if diskErr != nil {
		slog.Warn("Failed to collect metric", "metric", "disk", "error", diskErr)
		disk = 0
	}

	if uptimeErr != nil {
		slog.Warn("Failed to collect metric", "metric", "uptime", "error", uptimeErr)
		uptime = ""
	}

	if processesErr != nil {
		slog.Warn("Failed to collect metric", "metric", "processes", "error", processesErr)
		processes = 0
	}

	return Metrics{
		Hostname:  hostname,
		CPU:       cpu,
		Memory:    memory,
		Disk:      disk,
		Uptime:    uptime,
		Processes: processes,
	}
}
