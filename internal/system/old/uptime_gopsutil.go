package system

import (
	"context"

	"github.com/shirou/gopsutil/v4/host"
)

type GopsutilUptime struct{}

func (g GopsutilUptime) Uptime(ctx context.Context) (uint64, error) {
	return host.UptimeWithContext(ctx)
}
