package hardware

import (
	"context"
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
	IO          io.IOStatistics                    `json:"io"`
	Memory      memory.MemoryInformation           `json:"memory"`
	Motherboard motherboard.MotherboardInformation `json:"motherboard"`
	Network     []network.NetworkInformation       `json:"network"`
	Peripherals peripherals.PeripheralsInformation `json:"peripherals"`
	Sensors     []sensors.SensorInformation        `json:"sensors"`
	Storage     []storage.StorageInformation       `json:"storage"`
}

func CollectAllHardwareInformation(ctx context.Context) (CompleteHardwareInformation, error) {
	var info CompleteHardwareInformation
	var wg sync.WaitGroup
	var mu sync.Mutex

	logError := func(section string, err error) {
		slog.ErrorContext(ctx, "Failed to collect hardware section",
			slog.String("section", section),
			slog.String("error_details", err.Error()),
		)
	}

	// 1. Battery
	wg.Add(1)
	go func() {
		defer wg.Done()
		b := battery.HardwareBattery{}
		res, err := b.RetrieveBatteryInfo()
		if err != nil {
			logError("battery", err)
			return
		}
		mu.Lock()
		info.Battery = res
		mu.Unlock()
	}()

	// 2. BIOS
	wg.Add(1)
	go func() {
		defer wg.Done()
		b := bios.HardwareBios{}
		res, err := b.RetrieveBiosInformation()
		if err != nil {
			logError("bios", err)
			return
		}
		mu.Lock()
		info.BIOS = res
		mu.Unlock()
	}()

	// 3. CPU
	wg.Add(1)
	go func() {
		defer wg.Done()
		c := cpu.HardwareCpu{}
		res, err := c.RetrieveCPUInfo()
		if err != nil {
			logError("cpu", err)
			return
		}
		mu.Lock()
		info.CPU = res
		mu.Unlock()
	}()

	// 4. GPU
	wg.Add(1)
	go func() {
		defer wg.Done()
		g := gpu.HardwareGpu{}
		res, err := g.RetrieveGPUInfo()
		if err != nil {
			logError("gpu", err)
			return
		}
		mu.Lock()
		info.GPU = res
		mu.Unlock()
	}()

	// 5. IO
	wg.Add(1)
	go func() {
		defer wg.Done()
		i := io.HardwareIo{}
		res, err := i.RetrieveIOStats()
		if err != nil {
			logError("io", err)
			return
		}
		mu.Lock()
		info.IO = res
		mu.Unlock()
	}()

	// 6. Memory
	wg.Add(1)
	go func() {
		defer wg.Done()
		m := memory.HardwareMemory{}
		res, err := m.RetrieveMemoryInfo()
		if err != nil {
			logError("memory", err)
			return
		}
		mu.Lock()
		info.Memory = res
		mu.Unlock()
	}()

	// 7. Motherboard
	wg.Add(1)
	go func() {
		defer wg.Done()
		mb := motherboard.HardwareMotherboard{}
		res, err := mb.RetrieveMotherboardInfo()
		if err != nil {
			logError("motherboard", err)
			return
		}
		mu.Lock()
		info.Motherboard = res
		mu.Unlock()
	}()

	// 8. Network
	wg.Add(1)
	go func() {
		defer wg.Done()
		netImpl := network.HardwareNetwork{}
		res, err := netImpl.RetrieveNetworkInfo()
		if err != nil {
			logError("network", err)
			return
		}
		mu.Lock()
		info.Network = res
		mu.Unlock()
	}()

	// 9. Peripherals
	wg.Add(1)
	go func() {
		defer wg.Done()
		p := peripherals.HardwarePeripherals{}
		res, err := p.RetrievePeripheralsInfo()
		if err != nil {
			logError("peripherals", err)
			return
		}
		mu.Lock()
		info.Peripherals = res
		mu.Unlock()
	}()

	// 10. Sensors
	wg.Add(1)
	go func() {
		defer wg.Done()
		s := sensors.HardwareSensors{}
		res, err := s.RetrieveSensorsInfo()
		if err != nil {
			logError("sensors", err)
			return
		}
		mu.Lock()
		info.Sensors = res
		mu.Unlock()
	}()

	// 11. Storage
	wg.Add(1)
	go func() {
		defer wg.Done()
		st := storage.HardwareStorage{}
		res, err := st.RetrieveStorageInformation()
		if err != nil {
			logError("storage", err)
			return
		}
		mu.Lock()
		info.Storage = res
		mu.Unlock()
	}()

	wg.Wait()
	slog.InfoContext(ctx, "Successfully collected all hardware information")
	return info, nil
}
