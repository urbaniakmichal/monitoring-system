package testutil

import "context"

type FakeMemory struct {
	Value int
	Err   error
}

func (f FakeMemory) Usage(ctx context.Context) (int, error) {
	return f.Value, f.Err
}
