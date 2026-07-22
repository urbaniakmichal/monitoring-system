package testutil

import "context"

type FakeDisk struct {
	Value int
	Err   error
}

func (f FakeDisk) Usage(ctx context.Context, drive string) (int, error) {
	return f.Value, f.Err
}
