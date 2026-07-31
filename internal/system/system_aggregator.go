package system

import (
	"fmt"
	"log/slog"
	"sync"

	"monitoring-system/internal/system/load"
	system "monitoring-system/internal/system/os"
	"monitoring-system/internal/system/users"
)

type CompleteSystemInformation struct {
	System system.SystemInformation `json:"system"`
	Load   load.LoadInformation     `json:"load"`
	Users  []users.UserInformation  `json:"users"`
}

func CollectAllSystemInformation() (CompleteSystemInformation, error) {
	var info CompleteSystemInformation
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

	// 1. System general info (hostname, OS, uptime, top processes)
	wg.Add(1)
	go func() {
		defer wg.Done()
		res, err := system.RetrieveSystemInfo()
		if err != nil {
			slog.Error("Failed to collect system info during aggregation", slog.String("error_details", err.Error()))
			setError(fmt.Errorf("failed to collect system info: %w", err))
			return
		}
		mu.Lock()
		info.System = res
		mu.Unlock()
	}()

	// 2. System load averages (load1, load5, load15)
	wg.Add(1)
	go func() {
		defer wg.Done()
		res, err := load.RetrieveLoadInfo()
		if err != nil {
			slog.Error("Failed to collect load info during aggregation", slog.String("error_details", err.Error()))
			setError(fmt.Errorf("failed to collect load info: %w", err))
			return
		}
		mu.Lock()
		info.Load = res
		mu.Unlock()
	}()

	// 3. Active users info
	wg.Add(1)
	go func() {
		defer wg.Done()
		res, err := users.RetrieveUsersInfo()
		if err != nil {
			slog.Error("Failed to collect users info during aggregation", slog.String("error_details", err.Error()))
			setError(fmt.Errorf("failed to collect users info: %w", err))
			return
		}
		mu.Lock()
		info.Users = res
		mu.Unlock()
	}()

	// Wait for all goroutines to complete
	wg.Wait()

	if firstErr != nil {
		return CompleteSystemInformation{}, firstErr
	}

	slog.Info("Successfully collected all system metrics and environment info concurrently")
	return info, nil
}