//go:build linux || darwin

package storage

import (
	"fmt"
	"log/slog"
	"syscall"
)

func RetrieveStorageInfo(path string) (StorageInformation, error) {
	var stat syscall.Statfs_t
	err := syscall.Statfs(path, &stat)
	if err != nil {
		slog.Error("Failed to get storage stats", slog.String("path", path), slog.String("error_details", err.Error()))
		return StorageInformation{}, fmt.Errorf("failed to get storage info: %w", err)
	}

	totalBytes := stat.Blocks * uint64(stat.Bsize)
	freeBytes := stat.Bavail * uint64(stat.Bsize)
	totalMB := totalBytes / (1024 * 1024)
	freeMB := freeBytes / (1024 * 1024)
	usedMB := totalMB - freeMB

	var usedPercent float64
	if totalMB > 0 {
		usedPercent = (float64(usedMB) / float64(totalMB)) * 100
	}

	return StorageInformation{
		Path:        path,
		TotalMB:     totalMB,
		FreeMB:      freeMB,
		UsedPercent: usedPercent,
	}, nil
}