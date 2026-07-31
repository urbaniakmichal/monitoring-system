package metrics

import (
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// --- SYSTEM METRICS (3/3) ---
	load1Gauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "monitor_system_load_1",
			Help: "System load average over 1 minute",
		},
		[]string{"hostname"},
	)

	load5Gauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "monitor_system_load_5",
			Help: "System load average over 5 minutes",
		},
		[]string{"hostname"},
	)

	load15Gauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "monitor_system_load_15",
			Help: "System load average over 15 minutes",
		},
		[]string{"hostname"},
	)

	systemUsersCountGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "monitor_system_users_total",
			Help: "Total number of active or logged users detected",
		},
		[]string{"hostname"},
	)

	topProcessesGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "monitor_top_processes_tracked_total",
			Help: "Total number of tracked top processes",
		},
		[]string{"hostname"},
	)

	uptimeGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "monitor_uptime_seconds",
			Help: "Total uptime of the system in seconds",
		},
		[]string{"hostname"},
	)

	// --- HARDWARE METRICS (11/11) ---
	cpuUsageGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "monitor_cpu_usage_percent",
			Help: "Current CPU usage percentage",
		},
		[]string{"hostname"},
	)

	memoryUsageGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "monitor_memory_usage_percent",
			Help: "Current memory usage percentage",
		},
		[]string{"hostname"},
	)

	diskUsageGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "monitor_disk_usage_percent",
			Help: "Current disk usage percentage per device and path",
		},
		[]string{"hostname", "device", "path"},
	)

	hardwareBatteryCountGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "monitor_hardware_batteries_total",
			Help: "Total number of batteries detected",
		},
		[]string{"hostname"},
	)

	hardwareGpuCountGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "monitor_hardware_gpus_total",
			Help: "Total number of GPUs detected",
		},
		[]string{"hostname"},
	)

	hardwareNetworkCountGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "monitor_hardware_network_interfaces_total",
			Help: "Total number of hardware network interfaces detected",
		},
		[]string{"hostname"},
	)

	hardwareSensorsCountGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "monitor_hardware_sensors_total",
			Help: "Total number of hardware sensors detected",
		},
		[]string{"hostname"},
	)

	hardwareBiosInfoGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "monitor_hardware_bios_info",
			Help: "BIOS detailed information (value is always 1)",
		},
		[]string{"hostname"},
	)

	hardwareMotherboardInfoGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "monitor_hardware_motherboard_info",
			Help: "Motherboard detailed information (value is always 1)",
		},
		[]string{"hostname"},
	)

	hardwareIoInfoGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "monitor_hardware_io_info",
			Help: "IO statistics indicator (value is always 1)",
		},
		[]string{"hostname"},
	)

	hardwarePeripheralsInfoGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "monitor_hardware_peripherals_info",
			Help: "Peripherals indicator (value is always 1)",
		},
		[]string{"hostname"},
	)

	// --- SOFTWARE METRICS (9/9) ---
	softwareInfoGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "monitor_software_os_info",
			Help: "Operating system detailed information (value is always 1)",
		},
		[]string{"hostname"},
	)

	softwareUpdatesCountGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "monitor_software_updates_total",
			Help: "Total number of system updates detected",
		},
		[]string{"hostname"},
	)

	softwareDriversCountGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "monitor_software_drivers_total",
			Help: "Total number of system drivers detected",
		},
		[]string{"hostname"},
	)

	softwareNetworkAdaptersCountGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "monitor_software_network_adapters_total",
			Help: "Total number of network adapters detected",
		},
		[]string{"hostname"},
	)

	softwareApplicationsCountGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "monitor_software_applications_total",
			Help: "Total number of installed applications detected",
		},
		[]string{"hostname"},
	)

	softwareServicesCountGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "monitor_software_services_total",
			Help: "Total number of system services detected",
		},
		[]string{"hostname"},
	)

	softwareStartupCommandsCountGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "monitor_software_startup_commands_total",
			Help: "Total number of startup commands detected",
		},
		[]string{"hostname"},
	)

	softwareTasksCountGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "monitor_software_scheduled_tasks_total",
			Help: "Total number of scheduled tasks detected",
		},
		[]string{"hostname"},
	)

	softwareOptionalFeaturesCountGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "monitor_software_optional_features_total",
			Help: "Total number of optional features detected",
		},
		[]string{"hostname"},
	)
)

// RecordMetrics updates Prometheus gauges using the comprehensive nested Metrics struct.
func RecordMetrics(m Metrics) {
	hostname := m.System.System.Hostname
	if hostname == "" {
		hostname = "unknown"
	}

	// 1. System Metrics (3 fields)
	load1Gauge.WithLabelValues(hostname).Set(m.System.Load.Load1)
	load5Gauge.WithLabelValues(hostname).Set(m.System.Load.Load5)
	load15Gauge.WithLabelValues(hostname).Set(m.System.Load.Load15)
	systemUsersCountGauge.WithLabelValues(hostname).Set(float64(len(m.System.Users)))
	topProcessesGauge.WithLabelValues(hostname).Set(float64(len(m.System.System.TopProcesses)))

	uptimeStr := m.System.System.Uptime
	if strings.Contains(uptimeStr, "seconds") {
		fields := strings.Fields(uptimeStr)
		if len(fields) > 0 {
			if secs, err := strconv.ParseFloat(fields[0], 64); err == nil {
				uptimeGauge.WithLabelValues(hostname).Set(secs)
			}
		}
	} else if duration, err := time.ParseDuration(uptimeStr); err == nil {
		uptimeGauge.WithLabelValues(hostname).Set(duration.Seconds())
	}

	// 2. Hardware Metrics (11 fields)
	cpuUsageGauge.WithLabelValues(hostname).Set(m.Hardware.CPU.UsagePercent)
	memoryUsageGauge.WithLabelValues(hostname).Set(m.Hardware.Memory.UsedPercent)

	for _, disk := range m.Hardware.Storage {
		diskUsageGauge.WithLabelValues(hostname, disk.Device, disk.Path).Set(disk.UsedPercent)
	}

	hardwareBatteryCountGauge.WithLabelValues(hostname).Set(float64(len(m.Hardware.Battery)))
	hardwareGpuCountGauge.WithLabelValues(hostname).Set(float64(len(m.Hardware.GPU)))
	hardwareNetworkCountGauge.WithLabelValues(hostname).Set(float64(len(m.Hardware.Network)))
	hardwareSensorsCountGauge.WithLabelValues(hostname).Set(float64(len(m.Hardware.Sensors)))

	hardwareBiosInfoGauge.WithLabelValues(hostname).Set(1)
	hardwareMotherboardInfoGauge.WithLabelValues(hostname).Set(1)
	hardwareIoInfoGauge.WithLabelValues(hostname).Set(1)
	hardwarePeripheralsInfoGauge.WithLabelValues(hostname).Set(1)

	// 3. Software Metrics (9 fields)
	softwareInfoGauge.WithLabelValues(hostname).Set(1)
	softwareUpdatesCountGauge.WithLabelValues(hostname).Set(float64(len(m.Software.SystemUpdates)))
	softwareDriversCountGauge.WithLabelValues(hostname).Set(float64(len(m.Software.SystemDrivers)))
	softwareNetworkAdaptersCountGauge.WithLabelValues(hostname).Set(float64(len(m.Software.NetworkAdapters)))
	softwareApplicationsCountGauge.WithLabelValues(hostname).Set(float64(len(m.Software.InstalledApplications)))
	softwareServicesCountGauge.WithLabelValues(hostname).Set(float64(len(m.Software.SystemServices)))
	softwareStartupCommandsCountGauge.WithLabelValues(hostname).Set(float64(len(m.Software.StartupCommands)))
	softwareTasksCountGauge.WithLabelValues(hostname).Set(float64(len(m.Software.ScheduledTasks)))
}