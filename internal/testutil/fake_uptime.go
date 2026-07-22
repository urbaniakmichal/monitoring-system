package testutil

import "context"

type FakeUptime struct {
	Value uint64
	Err   error
}

func (f FakeUptime) Uptime(ctx context.Context) (uint64, error) {
	return f.Value, f.Err
}
