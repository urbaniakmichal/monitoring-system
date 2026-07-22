package system

import (
	"context"
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

	uptimeSeconds := uptime

	uptimeDuration := time.Duration(uptimeSeconds) * time.Second
	return uptimeDuration.String(), nil
}
