package system

import "context"

type MemoryProvider interface {
	Usage(ctx context.Context) (int, error)
}
