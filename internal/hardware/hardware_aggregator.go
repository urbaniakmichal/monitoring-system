package hardware

import (
	"fmt"
	"log/slog"
	"sync"

	"monitoring-system/internal/hardware/battery"
	"monitoring-system/internal/hardware/bios"
	"monitoring-system/internal/hardware/cpu"
	"monitoring-system/internal/hardware/gpu"
	"monitoring-system/internal/hardware/io"
	"monitoring-system/internal/hardware/memory"
	"monitoring-system/internal/hardware/motherboard"
	"monitoring-system/internal/hardware/network"
	"monitoring-system/internal/hardware/peripherals"
	"monitoring-system/internal/hardware/sensors"
	"monitoring-system/internal/hardware/storage"
)

type CompleteHardwareInformation struct {
	Battery     []battery.BatteryInformation       `json:"battery"`
	BIOS        bios.BiosInformation               `json:"bios"`
	CPU         cpu.CPUInformation                 `json:"cpu"`
	GPU         []gpu.GPUInformation               `json:"gpu"`
	IO          io.IOStatistics                 `json:"io"`
	Memory      memory.MemoryInformation           `json:"memory"`
	Motherboard motherboard.MotherboardInformation `json:"motherboard"`
	Network     []network.NetworkInformation  `json:"network"`
	Peripherals peripherals.PeripheralsInformation  `json:"peripherals"`
	Sensors     []sensors.SensorInformation          `json:"sensors"`
	Storage     []storage.StorageInformation       `json:"storage"`
}

func CollectAllHardwareInformation() (CompleteHardwareInformation, error) {
	var info CompleteHardwareInformation
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

	// 1. Battery information
	wg.Add(1)
	go func() {
		defer wg.Done()

		battery := battery.HardwareBattery{}
		res, err := battery.RetrieveBatteryInfo()

		if err != nil {
			slog.Error("Failed to collect battery info during aggregation", slog.String("error_details", err.Error()))
			setError(fmt.Errorf("failed to collect battery info: %w", err))
			return
		}
		mu.Lock()
		info.Battery = res
		mu.Unlock()
	}()

	// 2. BIOS information
	wg.Add(1)
	go func() {
		defer wg.Done()

		bios := bios.HardwareBios{}
		res, err := bios.RetrieveBiosInformation()

		if err != nil {
			slog.Error("Failed to collect BIOS info during aggregation", slog.String("error_details", err.Error()))
			setError(fmt.Errorf("failed to collect BIOS info: %w", err))
			return
		}
		mu.Lock()
		info.BIOS = res
		mu.Unlock()
	}()

	// 3. CPU information
	wg.Add(1)
	go func() {
		defer wg.Done()
		res, err := cpu.RetrieveCPUInfo()
		if err != nil {
			slog.Error("Failed to collect CPU info during aggregation", slog.String("error_details", err.Error()))
			setError(fmt.Errorf("failed to collect CPU info: %w", err))
			return
		}
		mu.Lock()
		info.CPU = res
		mu.Unlock()
	}()

	// 4. GPU information
	wg.Add(1)
	go func() {
		defer wg.Done()
		res, err := gpu.RetrieveGPUInfo()
		if err != nil {
			slog.Error("Failed to collect GPU info during aggregation", slog.String("error_details", err.Error()))
			setError(fmt.Errorf("failed to collect GPU info: %w", err))
			return
		}
		mu.Lock()
		info.GPU = res
		mu.Unlock()
	}()

	// 5. IO ports/devices information
	wg.Add(1)
	go func() {
		defer wg.Done()
		res, err := io.RetrieveIOStats()
		if err != nil {
			slog.Error("Failed to collect IO info during aggregation", slog.String("error_details", err.Error()))
			setError(fmt.Errorf("failed to collect IO info: %w", err))
			return
		}
		mu.Lock()
		info.IO = res
		mu.Unlock()
	}()

	// 6. Memory information
	wg.Add(1)
	go func() {
		defer wg.Done()
		res, err := memory.RetrieveMemoryInfo()
		if err != nil {
			slog.Error("Failed to collect memory info during aggregation", slog.String("error_details", err.Error()))
			setError(fmt.Errorf("failed to collect memory info: %w", err))
			return
		}
		mu.Lock()
		info.Memory = res
		mu.Unlock()
	}()

	// 7. Motherboard information
	wg.Add(1)
	go func() {
		defer wg.Done()
		res, err := motherboard.RetrieveMotherboardInfo()
		if err != nil {
			slog.Error("Failed to collect motherboard info during aggregation", slog.String("error_details", err.Error()))
			setError(fmt.Errorf("failed to collect motherboard info: %w", err))
			return
		}
		mu.Lock()
		info.Motherboard = res
		mu.Unlock()
	}()

	// 8. Network adapters information
	wg.Add(1)
	go func() {
		defer wg.Done()
		res, err := network.RetrieveNetworkInfo()
		if err != nil {
			slog.Error("Failed to collect network info during aggregation", slog.String("error_details", err.Error()))
			setError(fmt.Errorf("failed to collect network info: %w", err))
			return
		}
		mu.Lock()
		info.Network = res
		mu.Unlock()
	}()

	// 9. Peripherals information
	wg.Add(1)
	go func() {
		defer wg.Done()
		res, err := peripherals.RetrievePeripheralsInfo()
		if err != nil {
			slog.Error("Failed to collect peripherals info during aggregation", slog.String("error_details", err.Error()))
			setError(fmt.Errorf("failed to collect peripherals info: %w", err))
			return
		}
		mu.Lock()
		info.Peripherals = res
		mu.Unlock()
	}()

	// 10. Sensors information
	wg.Add(1)
	go func() {
		defer wg.Done()
		res, err := sensors.RetrieveSensorsInfo()
		if err != nil {
			slog.Error("Failed to collect sensors info during aggregation", slog.String("error_details", err.Error()))
			setError(fmt.Errorf("failed to collect sensors info: %w", err))
			return
		}
		mu.Lock()
		info.Sensors = res
		mu.Unlock()
	}()

	// 11. Storage information
	wg.Add(1)
	go func() {
		defer wg.Done()
		res, err := storage.RetrieveStorageInformation()
		if err != nil {
			slog.Error("Failed to collect storage info during aggregation", slog.String("error_details", err.Error()))
			setError(fmt.Errorf("failed to collect storage info: %w", err))
			return
		}
		mu.Lock()
		info.Storage = res
		mu.Unlock()
	}()

	// Wait for all goroutines to complete
	wg.Wait()

	if firstErr != nil {
		return CompleteHardwareInformation{}, firstErr
	}

	slog.Info("Successfully collected all hardware information concurrently")
	return info, nil
}