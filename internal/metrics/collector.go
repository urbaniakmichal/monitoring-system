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
	_, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var hwInfo hardware.CompleteHardwareInformation
	var swInfo software.CompleteSoftwareInformation
	var sysInfo system.CompleteSystemInformation
	var wg sync.WaitGroup
	var mu sync.Mutex
	var firstErr error

	setError := func(err error) {
		if err != nil {
			mu.Lock()
			if firstErr == nil {
				firstErr = err
			}
			mu.Unlock()
		}
	}

	// 1. Collect all hardware information concurrently
	wg.Add(1)
	go func() {
		defer wg.Done()
		res, err := hardware.CollectAllHardwareInformation()
		if err != nil {
			slog.Error("Failed to collect hardware information in metrics collector", slog.String("error", err.Error()))
			setError(fmt.Errorf("hardware collection error: %w", err))
			return
		}
		mu.Lock()
		hwInfo = res
		mu.Unlock()
	}()

	// 2. Collect all software information concurrently
	wg.Add(1)
	go func() {
		defer wg.Done()
		res, err := software.CollectAllSoftwareInformation()
		if err != nil {
			slog.Error("Failed to collect software information in metrics collector", slog.String("error", err.Error()))
			setError(fmt.Errorf("software collection error: %w", err))
			return
		}
		mu.Lock()
		swInfo = res
		mu.Unlock()
	}()

	// 3. Collect all system information concurrently
	wg.Add(1)
	go func() {
		defer wg.Done()
		res, err := system.CollectAllSystemInformation()
		if err != nil {
			slog.Error("Failed to collect system information in metrics collector", slog.String("error", err.Error()))
			setError(fmt.Errorf("system collection error: %w", err))
			return
		}
		mu.Lock()
		sysInfo = res
		mu.Unlock()
	}()

	// Wait for all hardware, software, and system aggregations to complete
	wg.Wait()

	if firstErr != nil {
		return Metrics{}, firstErr
	}

	return Metrics{
		Timestamp: time.Now(),
		Hardware:  hwInfo,
		Software:  swInfo,
		System:    sysInfo,
	}, nil
}