package testutil

import "context"

type FakeCPU struct {
	Value float64
	Err   error
}

func (f FakeCPU) Percent(ctx context.Context) ([]float64, error) {
	return []float64{f.Value}, f.Err
}
