package system

import (
	"context"

	"github.com/shirou/gopsutil/v4/disk"
)

type GopsutilDisk struct{}

func (g GopsutilDisk) Usage(ctx context.Context, drive string) (int, error) {

	stat, err := disk.UsageWithContext(ctx, drive)
	if err != nil {
		return 0, err
	}

	return int(stat.UsedPercent), nil
}
