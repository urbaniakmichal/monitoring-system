package software

import (
	"fmt"
	"log/slog"
	"monitoring-system/internal/software/applications"
	"monitoring-system/internal/software/drivers"
	"monitoring-system/internal/software/features"
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
	OptionalFeatures      []features.OptionalFeatureInformation
}

func CollectAllSoftwareInformation() (CompleteSoftwareInformation, error) {
	operatingSystemInformation, operatingSystemError := operatingSystem.RetrieveOperatingSystemInformation()
	if operatingSystemError != nil {
		slog.Error("Failed to collect operating system information during aggregation",
			slog.String("error_details", operatingSystemError.Error()),
		)
		return CompleteSoftwareInformation{}, fmt.Errorf("failed to collect operating system info: %w", operatingSystemError)
	}

	systemUpdatesList, systemUpdatesError := updates.RetrieveSystemUpdates()
	if systemUpdatesError != nil {
		slog.Error("Failed to collect system updates during aggregation",
			slog.String("error_details", systemUpdatesError.Error()),
		)
		return CompleteSoftwareInformation{}, fmt.Errorf("failed to collect system updates: %w", systemUpdatesError)
	}

	systemDriversList, systemDriversError := drivers.RetrieveInstalledDrivers()
	if systemDriversError != nil {
		slog.Error("Failed to collect system drivers during aggregation",
			slog.String("error_details", systemDriversError.Error()),
		)
		return CompleteSoftwareInformation{}, fmt.Errorf("failed to collect system drivers: %w", systemDriversError)
	}

	networkAdaptersList, networkAdaptersError := network.RetrieveActiveNetworkAdapters()
	if networkAdaptersError != nil {
		slog.Error("Failed to collect network adapters during aggregation",
			slog.String("error_details", networkAdaptersError.Error()),
		)
		return CompleteSoftwareInformation{}, fmt.Errorf("failed to collect network adapters: %w", networkAdaptersError)
	}

	installedApplicationsList, installedApplicationsError := applications.RetrieveInstalledApplications()
	if installedApplicationsError != nil {
		slog.Error("Failed to collect installed applications during aggregation",
			slog.String("error_details", installedApplicationsError.Error()),
		)
		return CompleteSoftwareInformation{}, fmt.Errorf("failed to collect installed applications: %w", installedApplicationsError)
	}

	systemServicesList, systemServicesError := services.RetrieveSystemServices()
	if systemServicesError != nil {
		slog.Error("Failed to collect system services during aggregation",
			slog.String("error_details", systemServicesError.Error()),
		)
		return CompleteSoftwareInformation{}, fmt.Errorf("failed to collect system services: %w", systemServicesError)
	}

	startupCommandsList, startupCommandsError := startup.RetrieveStartupCommands()
	if startupCommandsError != nil {
		slog.Error("Failed to collect startup commands during aggregation",
			slog.String("error_details", startupCommandsError.Error()),
		)
		return CompleteSoftwareInformation{}, fmt.Errorf("failed to collect startup commands: %w", startupCommandsError)
	}

	scheduledTasksList, scheduledTasksError := tasks.RetrieveScheduledTasks()
	if scheduledTasksError != nil {
		slog.Error("Failed to collect scheduled tasks during aggregation",
			slog.String("error_details", scheduledTasksError.Error()),
		)
		return CompleteSoftwareInformation{}, fmt.Errorf("failed to collect scheduled tasks: %w", scheduledTasksError)
	}

	optionalFeaturesList, optionalFeaturesError := features.RetrieveOptionalFeatures()
	if optionalFeaturesError != nil {
		slog.Error("Failed to collect optional features during aggregation",
			slog.String("error_details", optionalFeaturesError.Error()),
		)
		return CompleteSoftwareInformation{}, fmt.Errorf("failed to collect optional features: %w", optionalFeaturesError)
	}

	completeInformation := CompleteSoftwareInformation{
		OperatingSystem:       operatingSystemInformation,
		SystemUpdates:         systemUpdatesList,
		SystemDrivers:         systemDriversList,
		NetworkAdapters:       networkAdaptersList,
		InstalledApplications: installedApplicationsList,
		SystemServices:        systemServicesList,
		StartupCommands:       startupCommandsList,
		ScheduledTasks:        scheduledTasksList,
		OptionalFeatures:      optionalFeaturesList,
	}

	slog.Info("Successfully collected all software and system environment information")
	return completeInformation, nil
}
