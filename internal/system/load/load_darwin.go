//go:build darwin

package load

import (
	"os/exec"
	"strconv"
	"strings"
)

func (*Load)RetrieveLoadInfo() (LoadInformation, error) {
	cmd := exec.Command("sysctl", "-n", "vm.loadavg")
	output, err := cmd.Output()
	if err != nil {
		return LoadInformation{}, nil
	}
	cleaned := strings.Trim(string(output), "{} \n")
	parts := strings.Fields(cleaned)
	if len(parts) < 3 {
		return LoadInformation{}, nil
	}
	l1, _ := strconv.ParseFloat(parts[0], 64)
	l5, _ := strconv.ParseFloat(parts[1], 64)
	l15, _ := strconv.ParseFloat(parts[2], 64)

	return LoadInformation{Load1: l1, Load5: l5, Load15: l15}, nil
}