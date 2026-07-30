package system

import "context"

type CPUProvider interface {
	Percent(ctx context.Context) ([]float64, error)
}
