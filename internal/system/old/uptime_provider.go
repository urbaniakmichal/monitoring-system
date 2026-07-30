package system

import "context"

type UptimeProvider interface {
	Uptime(ctx context.Context) (uint64, error)
}
