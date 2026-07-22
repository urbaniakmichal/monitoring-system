package platform

import "runtime"

type Platform string

const (
	Windows Platform = "windows"
	Linux   Platform = "linux"
	MacOS   Platform = "darwin"
	Unknown Platform = "unknown"
)

func Current() Platform {

	switch runtime.GOOS {
	case "windows":
		return Windows

	case "linux":
		return Linux

	case "darwin":
		return MacOS

	default:
		return Unknown
	}
}
