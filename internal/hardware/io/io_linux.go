//go:build linux

package io

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func (*HardwareIo)RetrieveIOStats() (IOStatistics, error) {
	netStats := make(map[string]NetworkIO)
	if file, err := os.Open("/proc/net/dev"); err == nil {
		defer func() { _ = file.Close() }()
		scanner := bufio.NewScanner(file)
		i := 0
		for scanner.Scan() {
			i++
			if i <= 2 {
				continue
			}
			line := strings.TrimSpace(scanner.Text())
			parts := strings.Split(line, ":")
			if len(parts) != 2 {
				continue
			}
			iface := strings.TrimSpace(parts[0])
			fields := strings.Fields(parts[1])
			if len(fields) >= 9 {
				recv, _ := strconv.ParseUint(fields[0], 10, 64)
				sent, _ := strconv.ParseUint(fields[8], 10, 64)
				netStats[iface] = NetworkIO{BytesSent: sent, BytesRecv: recv}
			}
		}
	}

	diskStats := make(map[string]DiskIO)
	if file, err := os.Open("/proc/diskstats"); err == nil {
		defer func() { _ = file.Close() }()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			fields := strings.Fields(scanner.Text())
			if len(fields) >= 14 {
				devName := fields[2]
				sectorsRead, _ := strconv.ParseUint(fields[5], 10, 64)
				sectorsWritten, _ := strconv.ParseUint(fields[9], 10, 64)
				diskStats[devName] = DiskIO{
					ReadBytes:  sectorsRead * 512,
					WriteBytes: sectorsWritten * 512,
				}
			}
		}
	}

	return IOStatistics{Network: netStats, Disk: diskStats}, nil
}