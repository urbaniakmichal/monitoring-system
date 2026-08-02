//go:build linux

package storage

import (
	"bufio"
	"os"
	"strings"
	"syscall"
)

func RetrieveStorageInformation() ([]StorageInformation, error) {
	file, err := os.Open("/proc/mounts")
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	var storages []StorageInformation
	seenPaths := make(map[string]bool)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 3 {
			continue
		}
		device := fields[0]
		mountPath := fields[1]
		fstype := fields[2]

		if !isRealFilesystem(fstype) || seenPaths[mountPath] {
			continue
		}

		var stat syscall.Statfs_t
		if err := syscall.Statfs(mountPath, &stat); err != nil {
			continue
		}

		// #nosec G115
		totalBytes := uint64(stat.Blocks) * uint64(stat.Bsize)
		// #nosec G115
		freeBytes := uint64(stat.Bavail) * uint64(stat.Bsize)
		if totalBytes == 0 {
			continue
		}

		totalMB := totalBytes / (1024 * 1024)
		freeMB := freeBytes / (1024 * 1024)
		usedMB := totalMB - freeMB
		var usedPercent float64
		if totalMB > 0 {
			usedPercent = (float64(usedMB) / float64(totalMB)) * 100
		}

		seenPaths[mountPath] = true
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

func isRealFilesystem(fstype string) bool {
	allowed := map[string]bool{
		"ext3": true, "ext4": true, "xfs": true, "btrfs": true,
		"vfat": true, "ntfs": true, "zfs": true, "f2fs": true,
	}
	return allowed[fstype]
}