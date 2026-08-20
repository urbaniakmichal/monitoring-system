package system

import (
	"context"
	"log/slog"
	"sync"

	"monitoring-system/internal/system/load"
	systemOS "monitoring-system/internal/system/os"
	"monitoring-system/internal/system/users"
)

type CompleteSystemInformation struct {
	System systemOS.SystemInformation `json:"system"`
	Load   load.LoadInformation       `json:"load"`
	Users  []users.UserInformation    `json:"users"`
}

func CollectAllSystemInformation(ctx context.Context) (CompleteSystemInformation, error) {
	var info CompleteSystemInformation
	var wg sync.WaitGroup
	var mu sync.Mutex

	logError := func(section string, err error) {
		slog.ErrorContext(ctx, "Failed to collect system section",
			slog.String("section", section),
			slog.String("error_details", err.Error()),
		)
	}

	// 1. System general info
	wg.Add(1)
	go func() {
		defer wg.Done()
		sysImpl := systemOS.SystemOsInfo{}
		res, err := sysImpl.RetrieveSystemInfo()
		if err != nil {
			logError("system_os", err)
			return
		}
		mu.Lock()
		info.System = res
		mu.Unlock()
	}()

	// 2. System load averages
	wg.Add(1)
	go func() {
		defer wg.Done()
		loadImpl := load.Load{}
		res, err := loadImpl.RetrieveLoadInfo()
		if err != nil {
			logError("system_load", err)
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
		usersImpl := users.SystemUsers{}
		res, err := usersImpl.RetrieveUsersInfo()
		if err != nil {
			logError("system_users", err)
			return
		}
		mu.Lock()
		info.Users = res
		mu.Unlock()
	}()

	wg.Wait()
	slog.InfoContext(ctx, "Successfully collected all system metrics")
	return info, nil
}
