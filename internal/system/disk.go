package system

import (
	"context"
	"monitor-agent/internal/platform"
)

func Disk(ctx context.Context, provider DiskProvider) (int, error) {
	usedPercent, err := provider.Usage(ctx, platform.DefaultDisk())
	if err != nil {
		return 0, err
	}
	return usedPercent, nil
}
