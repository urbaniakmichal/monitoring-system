package system

import "context"

type ProcessesProvider interface {
	Pids(ctx context.Context) ([]int32, error)
}
