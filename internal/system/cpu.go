package system

import (
	"context"
	"math"
)

func CPU(ctx context.Context, provider CPUProvider) (int, error) {
	percent, err := provider.Percent(ctx)

	if err != nil || len(percent) == 0 {
		return 0, err
	}

	return int(math.Round(percent[0])), nil
}
