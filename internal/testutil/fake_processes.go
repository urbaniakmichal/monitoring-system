package testutil

import (
	"context"
)

type FakeProcesses struct {
	Value []int32
	Err   error
}

func (f FakeProcesses) Pids(ctx context.Context) ([]int32, error) {
	return f.Value, f.Err
}
