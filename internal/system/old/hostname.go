package system

import (
	"context"
)

func Hostname(ctx context.Context, provider HostnameProvider) (string, error) {
	hostname, err := provider.Hostname(ctx)
	if err != nil {
		return "", err
	}
	return hostname, nil
}
