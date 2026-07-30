//go:build linux

package load

import (
	"os"
	"strconv"
	"strings"
)

func RetrieveLoadInfo() (LoadInformation, error) {
	data, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return LoadInformation{}, err
	}
	fields := strings.Fields(string(data))
	if len(fields) < 3 {
		return LoadInformation{}, nil
	}
	l1, _ := strconv.ParseFloat(fields[0], 64)
	l5, _ := strconv.ParseFloat(fields[1], 64)
	l15, _ := strconv.ParseFloat(fields[2], 64)

	return LoadInformation{Load1: l1, Load5: l5, Load15: l15}, nil
}