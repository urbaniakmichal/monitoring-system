//go:build darwin

package storage

import (
	"bytes"
	"os/exec"
	"strconv"
	"strings"
)

func (*HardwareStorage)RetrieveStorageInformation() ([]StorageInformation, error) {
	cmd := exec.Command("df", "-k")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	var storages []StorageInformation
	lines := strings.Split(out.String(), "\n")
	for i, line := range lines {
		if i == 0 || strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 9 {
			continue
		}
		device := fields[0]
		totalKB, _ := strconv.ParseUint(fields[1], 10, 64)
		usedKB, _ := strconv.ParseUint(fields[2], 10, 64)
		freeKB, _ := strconv.ParseUint(fields[3], 10, 64)
		mountPath := fields[8]

		totalMB := totalKB / 1024
		freeMB := freeKB / 1024
		usedMB := usedKB / 1024
		var usedPercent float64
		if totalMB > 0 {
			usedPercent = (float64(usedMB) / float64(totalMB)) * 100
		}

		storages = append(storages, StorageInformation{
			Device:      device,
			Path:        mountPath,
			TotalMB:     totalMB,
			FreeMB:      freeMB,
			UsedPercent: usedPercent,
		})
	}
	return storages, nil
}