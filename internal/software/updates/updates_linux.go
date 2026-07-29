//go:build linux

package updates

import (
	"bytes"
	"log/slog"
	"os/exec"
	"strings"
)

type SystemUpdateInformation struct {
	HotFixID    string `json:"HotFixID"`
	Description string `json:"Description"`
	InstalledOn string `json:"InstalledOn"`
	InstalledBy string `json:"InstalledBy"`
}

func RetrieveSystemUpdates() ([]SystemUpdateInformation, error) {
	cmd := exec.Command("apt", "list", "--installed")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		slog.Info("Apt list not available or failed")
		return []SystemUpdateInformation{}, nil
	}

	var updates []SystemUpdateInformation
	lines := strings.Split(out.String(), "\n")
	for i, line := range lines {
		if i == 0 || strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Split(line, "/")
		if len(parts) > 0 {
			updates = append(updates, SystemUpdateInformation{
				HotFixID:    parts[0],
				Description: parts[0],
				InstalledOn: "",
				InstalledBy: "APT",
			})
		}
	}
	return updates, nil
}