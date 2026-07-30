package system

import (
	"context"
	"os"
)

type GopsutilHostname struct{}

func (g GopsutilHostname) Hostname(ctx context.Context) (string, error) {
	return os.Hostname()
}
