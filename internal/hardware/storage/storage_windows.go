//go:build windows

package storage

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"os/exec"
)

type winLogicalDisk struct {
	DeviceID  string  `json:"DeviceID"`
	Size      *uint64 `json:"Size"`
	FreeSpace *uint64 `json:"FreeSpace"`
}

func (*HardwareStorage)RetrieveStorageInformation() ([]StorageInformation, error) {
	cmd := exec.Command("powershell", "-NoProfile", "-Command", "Get-CimInstance Win32_LogicalDisk | Select-Object DeviceID, Size, FreeSpace | ConvertTo-Json")
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		slog.Error("Failed to get storage on Windows", slog.String("error_details", err.Error()))
		return nil, err
	}

	trimmed := bytes.TrimSpace(out.Bytes())
	if len(trimmed) == 0 {
		return []StorageInformation{}, nil
	}

	var disks []winLogicalDisk
	if trimmed[0] == '{' {
		var single winLogicalDisk
		if err := json.Unmarshal(trimmed, &single); err != nil {
			return nil, err
		}
		disks = []winLogicalDisk{single}
	} else {
		if err := json.Unmarshal(trimmed, &disks); err != nil {
			return nil, err
		}
	}

	var storages []StorageInformation
	for _, d := range disks {
		if d.Size == nil || *d.Size == 0 {
			continue
		}
		totalBytes := *d.Size
		var freeBytes uint64
		if d.FreeSpace != nil {
			freeBytes = *d.FreeSpace
		}

		totalMB := totalBytes / (1024 * 1024)
		freeMB := freeBytes / (1024 * 1024)
		usedMB := totalMB - freeMB
		var usedPercent float64
		if totalMB > 0 {
			usedPercent = (float64(usedMB) / float64(totalMB)) * 100
		}

		storages = append(storages, StorageInformation{
			Device:      d.DeviceID,
			Path:        d.DeviceID + "\\",
			TotalMB:     totalMB,
			FreeMB:      freeMB,
			UsedPercent: usedPercent,
		})
	}

	return storages, nil
}