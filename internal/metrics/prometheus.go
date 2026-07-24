package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	cpuUsageGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "monitor_cpu_usage_percent",
			Help: "Current CPU usage percentage",
		},
		[]string{"hostname"},
	)

	memoryGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "monitor_memory_usage_bytes",
			Help: "Current memory usage in bytes",
		},
		[]string{"hostname"},
	)

	diskUsageGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "monitor_disk_usage_percent",
			Help: "Current disk usage percentage",
		},
		[]string{"hostname"},
	)

	processesGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "monitor_processes_total",
			Help: "Total number of running processes",
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
)

// RecordMetrics updates Prometheus gauges using the application Metrics struct.
func RecordMetrics(m Metrics) {
	cpuUsageGauge.WithLabelValues(m.Hostname).Set(float64(m.CPU))
	memoryGauge.WithLabelValues(m.Hostname).Set(float64(m.Memory))
	diskUsageGauge.WithLabelValues(m.Hostname).Set(float64(m.Disk))
	processesGauge.WithLabelValues(m.Hostname).Set(float64(m.Processes))

	if duration, err := time.ParseDuration(m.Uptime); err == nil {
		uptimeGauge.WithLabelValues(m.Hostname).Set(duration.Seconds())
	}
}
