package system

import "context"

type HostnameProvider interface {
	Hostname(ctx context.Context) (string, error)
}
