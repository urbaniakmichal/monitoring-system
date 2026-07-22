package system

import (
	"context"
)

func Processes(ctx context.Context, provider ProcessesProvider) (int, error) {

	pids, err := provider.Pids(ctx)
	if err != nil {
		return 0, err
	}

	return len(pids), nil
}
