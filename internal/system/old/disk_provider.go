package system

import "context"

type DiskProvider interface {
	Usage(ctx context.Context, drive string) (int, error)
}
