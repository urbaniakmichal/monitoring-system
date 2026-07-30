package system

import (
	"context"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
)

type GopsutilCPU struct{}

func (g GopsutilCPU) Percent(ctx context.Context) ([]float64, error) {
	return cpu.PercentWithContext(
		ctx,
		200*time.Millisecond,
		false,
	)
}
