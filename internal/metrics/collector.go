package metrics

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"monitoring-system/internal/hardware"
	"monitoring-system/internal/software"
	"monitoring-system/internal/system"
)

// Collect gathers complete system metrics by running hardware, software, and system aggregators concurrently.
func Collect(ctx context.Context, timeout time.Duration) (Metrics, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var hwInfo hardware.CompleteHardwareInformation
	var swInfo software.CompleteSoftwareInformation
	var sysInfo system.CompleteSystemInformation

	var wg sync.WaitGroup

	// 1. Hardware - pass context with Trace ID
	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		res, err := hardware.CollectAllHardwareInformation(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "Hardware aggregation returned error", slog.String("error", err.Error()))
			return
		}
		hwInfo = res
	}(ctx)

	// 2. Software - pass context with Trace ID
	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		res, err := software.CollectAllSoftwareInformation(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "Software aggregation returned error", slog.String("error", err.Error()))
			return
		}
		swInfo = res
	}(ctx)

	// 3. System - pass context with Trace ID
	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		res, err := system.CollectAllSystemInformation(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "System aggregation returned error", slog.String("error", err.Error()))
			return
		}
		sysInfo = res
	}(ctx)

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		slog.InfoContext(ctx, "Successfully collected all metrics within timeout limit")
		return Metrics{
			Timestamp: time.Now().UTC(),
			Hardware:  hwInfo,
			Software:  swInfo,
			System:    sysInfo,
		}, nil

	case <-ctx.Done():
		slog.ErrorContext(ctx, "Metrics collection timed out", slog.Duration("timeout", timeout))
		return Metrics{}, fmt.Errorf("metrics collection timed out after %s: %w", timeout, ctx.Err())
	}
}
