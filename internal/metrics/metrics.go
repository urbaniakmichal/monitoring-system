package metrics

import (
	"time"

	"monitoring-system/internal/hardware"
	"monitoring-system/internal/software"
	"monitoring-system/internal/system"
)

type Metrics struct {
	TraceID   string                                   `json:"trace_id"`
	Timestamp time.Time                                `json:"timestamp"`
	Hardware  hardware.CompleteHardwareInformation     `json:"hardware"`
	Software  software.CompleteSoftwareInformation     `json:"software"`
	System    system.CompleteSystemInformation         `json:"system"`
}