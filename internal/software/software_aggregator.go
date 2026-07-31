package software

import (
	"fmt"
	"log/slog"
	"sync"

	"monitoring-system/internal/software/applications"
	"monitoring-system/internal/software/drivers"
	"monitoring-system/internal/software/network"
	operatingSystem "monitoring-system/internal/software/operating_system"
	"monitoring-system/internal/software/services"
	"monitoring-system/internal/software/startup"
	"monitoring-system/internal/software/tasks"
	"monitoring-system/internal/software/updates"
)

type CompleteSoftwareInformation struct {
	OperatingSystem       operatingSystem.OperatingSystemInformation
	SystemUpdates         []updates.SystemUpdateInformation
	SystemDrivers         []drivers.DriverInformation
	NetworkAdapters       []network.NetworkAdapterInformation
	InstalledApplications []applications.ApplicationInformation
	SystemServices        []services.ServiceInformation
	StartupCommands       []startup.StartupCommandInformation
	ScheduledTasks        []tasks.ScheduledTaskInformation
}

func CollectAllSoftwareInformation() (CompleteSoftwareInformation, error) {
	var info CompleteSoftwareInformation
	var wg sync.WaitGroup
	var mu sync.Mutex
	var firstErr error

	setError := func(err error) {
		if err != nil {
			mu.Lock()
			if firstErr == nil {
				firstErr = err
			}
			mu.Unlock()
		}
	}

	// 1. Operating system
	wg.Add(1)
	go func() {
		defer wg.Done()
		res, err := operatingSystem.RetrieveOperatingSystemInformation()
		if err != nil {
			slog.Error("Failed to collect operating system info during aggregation", slog.String("error_details", err.Error()))
			setError(fmt.Errorf("failed to collect operating system info: %w", err))
			return
		}
		mu.Lock()
		info.OperatingSystem = res
		mu.Unlock()
	}()

	// 2. System updates
	wg.Add(1)
	go func() {
		defer wg.Done()
		res, err := updates.RetrieveSystemUpdates()
		if err != nil {
			slog.Error("Failed to collect system updates during aggregation", slog.String("error_details", err.Error()))
			setError(fmt.Errorf("failed to collect system updates: %w", err))
			return
		}
		mu.Lock()
		info.SystemUpdates = res
		mu.Unlock()
	}()

	// 3. System drivers
	wg.Add(1)
	go func() {
		defer wg.Done()
		res, err := drivers.RetrieveInstalledDrivers()
		if err != nil {
			slog.Error("Failed to collect system drivers during aggregation", slog.String("error_details", err.Error()))
			setError(fmt.Errorf("failed to collect system drivers: %w", err))
			return
		}
		mu.Lock()
		info.SystemDrivers = res
		mu.Unlock()
	}()

	// 4. Network adapters
	wg.Add(1)
	go func() {
		defer wg.Done()
		res, err := network.RetrieveActiveNetworkAdapters()
		if err != nil {
			slog.Error("Failed to collect network adapters during aggregation", slog.String("error_details", err.Error()))
			setError(fmt.Errorf("failed to collect network adapters: %w", err))
			return
		}
		mu.Lock()
		info.NetworkAdapters = res
		mu.Unlock()
	}()

	// 5. Installed applications
	wg.Add(1)
	go func() {
		defer wg.Done()
		res, err := applications.RetrieveInstalledApplications()
		if err != nil {
			slog.Error("Failed to collect installed applications during aggregation", slog.String("error_details", err.Error()))
			setError(fmt.Errorf("failed to collect installed applications: %w", err))
			return
		}
		mu.Lock()
		info.InstalledApplications = res
		mu.Unlock()
	}()

	// 6. System services
	wg.Add(1)
	go func() {
		defer wg.Done()
		res, err := services.RetrieveSystemServices()
		if err != nil {
			slog.Error("Failed to collect system services during aggregation", slog.String("error_details", err.Error()))
			setError(fmt.Errorf("failed to collect system services: %w", err))
			return
		}
		mu.Lock()
		info.SystemServices = res
		mu.Unlock()
	}()

	// 7. Startup commands
	wg.Add(1)
	go func() {
		defer wg.Done()
		res, err := startup.RetrieveStartupCommands()
		if err != nil {
			slog.Error("Failed to collect startup commands during aggregation", slog.String("error_details", err.Error()))
			setError(fmt.Errorf("failed to collect startup commands: %w", err))
			return
		}
		mu.Lock()
		info.StartupCommands = res
		mu.Unlock()
	}()

	// 8. Scheduled tasks
	wg.Add(1)
	go func() {
		defer wg.Done()
		res, err := tasks.RetrieveScheduledTasks()
		if err != nil {
			slog.Error("Failed to collect scheduled tasks during aggregation", slog.String("error_details", err.Error()))
			setError(fmt.Errorf("failed to collect scheduled tasks: %w", err))
			return
		}
		mu.Lock()
		info.ScheduledTasks = res
		mu.Unlock()
	}()

	// Wait for all goroutines to complete
	wg.Wait()

	if firstErr != nil {
		return CompleteSoftwareInformation{}, firstErr
	}

	slog.Info("Successfully collected all software and system environment information concurrently")
	return info, nil
}