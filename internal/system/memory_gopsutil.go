package system

import (
	"context"

	"github.com/shirou/gopsutil/v4/mem"
)

type GopsutilMemory struct{}

func (g GopsutilMemory) Usage(ctx context.Context) (int, error) {

	data, err := mem.VirtualMemoryWithContext(ctx)

	if err != nil {
		return 0, err
	}

	return int(data.UsedPercent), nil
}
