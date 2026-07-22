package platform

import "runtime"

func DefaultDisk() string {

	switch runtime.GOOS {
	case "windows":
		return "C:"

	case "linux", "darwin":
		return "/"

	default:
		return "."
	}
}
