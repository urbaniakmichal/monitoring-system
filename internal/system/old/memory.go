package system

import (
	"context"
)

func Memory(ctx context.Context, provider MemoryProvider) (int, error) {

	usage, err := provider.Usage(ctx)

	if err != nil {
		return 0, err
	}

	return usage, nil
}
