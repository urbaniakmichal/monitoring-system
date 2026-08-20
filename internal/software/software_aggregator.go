package software

import (
	"context"
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
	OperatingSystem       operatingSystem.OperatingSystemInformation `json:"operating_system"`
	SystemUpdates         []updates.SystemUpdateInformation          `json:"system_updates"`
	SystemDrivers         []drivers.DriverInformation                `json:"system_drivers"`
	NetworkAdapters       []network.NetworkAdapterInformation        `json:"network_adapters"`
	InstalledApplications []applications.ApplicationInformation      `json:"installed_applications"`
	SystemServices        []services.ServiceInformation              `json:"system_services"`
	StartupCommands       []startup.StartupCommandInformation        `json:"startup_commands"`
	ScheduledTasks        []tasks.ScheduledTaskInformation           `json:"scheduled_tasks"`
}

func CollectAllSoftwareInformation(ctx context.Context) (CompleteSoftwareInformation, error) {
	var info CompleteSoftwareInformation
	var wg sync.WaitGroup
	var mu sync.Mutex

	logError := func(section string, err error) {
		slog.ErrorContext(ctx, "Failed to collect software section",
			slog.String("section", section),
			slog.String("error_details", err.Error()),
		)
	}

	// 1. Operating system
	wg.Add(1)
	go func() {
		defer wg.Done()
		osImpl := &operatingSystem.SoftwareOperatingSystem{}
		res, err := osImpl.RetrieveOperatingSystemInformation()
		if err != nil {
			logError("operating_system", err)
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
		updatesImpl := &updates.SoftwareUpdates{}
		res, err := updatesImpl.RetrieveSystemUpdates()
		if err != nil {
			logError("system_updates", err)
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
		driversImpl := &drivers.SoftwareDrivers{}
		res, err := driversImpl.RetrieveInstalledDrivers()
		if err != nil {
			logError("system_drivers", err)
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
		networksImpl := &network.SoftwareNetworks{}
		res, err := networksImpl.RetrieveActiveNetworkAdapters()
		if err != nil {
			logError("network_adapters", err)
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
		appsImpl := &applications.SoftwareApplications{}
		res, err := appsImpl.RetrieveInstalledApplications()
		if err != nil {
			logError("installed_applications", err)
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
		servicesImpl := &services.SoftwareServices{}
		res, err := servicesImpl.RetrieveSystemServices()
		if err != nil {
			logError("system_services", err)
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
		startupImpl := &startup.SoftwareStartup{}
		res, err := startupImpl.RetrieveStartupCommands()
		if err != nil {
			logError("startup_commands", err)
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
		tasksImpl := &tasks.SoftwareTasks{}
		res, err := tasksImpl.RetrieveScheduledTasks()
		if err != nil {
			logError("scheduled_tasks", err)
			return
		}
		mu.Lock()
		info.ScheduledTasks = res
		mu.Unlock()
	}()

	wg.Wait()
	slog.InfoContext(ctx, "Successfully collected software information")
	return info, nil
}
