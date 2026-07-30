package system

import (
	"context"
	"fmt"
	"math"
	"time"
)

func Uptime(ctx context.Context, provider UptimeProvider) (string, error) {
	uptime, err := provider.Uptime(ctx)
	if err != nil {
		return "", err
	}
	if uptime == 0 {
		return "0s", nil
	}

	maxUptimeSeconds := uint64(math.MaxInt64 / int64(time.Second))
	if uptime > maxUptimeSeconds {
		return "", fmt.Errorf("uptime value %d exceeds maximum supported duration", uptime)
	}

	//nolint:gosec // G115: safe conversion guaranteed by bounds check above
	uptimeDuration := time.Duration(uptime) * time.Second
	return uptimeDuration.String(), nil
}
