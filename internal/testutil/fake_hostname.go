package testutil

import "context"

type FakeHostname struct {
	Value string
	Err   error
}

func (f FakeHostname) Hostname(ctx context.Context) (string, error) {
	return f.Value, f.Err
}
