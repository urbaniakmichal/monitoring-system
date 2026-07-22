package system

import (
	"context"

	"github.com/shirou/gopsutil/v4/process"
)

type GopsutilProcesses struct{}

func (g GopsutilProcesses) Pids(ctx context.Context) ([]int32, error) {
	processes, err := process.ProcessesWithContext(ctx)

	if err != nil {
		return nil, err
	}

	result := make([]int32, 0, len(processes))

	for _, p := range processes {
		result = append(result, p.Pid)
	}

	return result, nil
}
