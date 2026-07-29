//go:build linux

package startup

import (
	"os"
	"path/filepath"
	"strings"
)

type StartupCommandInformation struct {
	Name     string `json:"Name"`
	Command  string `json:"Command"`
	Location string `json:"Location"`
	User     string `json:"User"`
}

func RetrieveStartupCommands() ([]StartupCommandInformation, error) {
	autostartDir := "/etc/xdg/autostart"
	var startups []StartupCommandInformation

	files, err := os.ReadDir(autostartDir)
	if err != nil {
		return []StartupCommandInformation{}, nil
	}

	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".desktop") {
			startups = append(startups, StartupCommandInformation{
				Name:     f.Name(),
				Command:  filepath.Join(autostartDir, f.Name()),
				Location: autostartDir,
				User:     "system",
			})
		}
	}
	return startups, nil
}