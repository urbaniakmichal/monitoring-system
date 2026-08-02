package operating_system

type OperatingSystemInformation struct {
	Caption        string `json:"Caption"`
	Version        string `json:"Version"`
	BuildNumber    string `json:"BuildNumber"`
	Manufacturer   string `json:"Manufacturer"`
	OSArchitecture string `json:"OSArchitecture"`
	InstallDate    string `json:"InstallDate"`
	LastBootTime   string `json:"LastBootUpTime"`
}

type OperatingSystem interface {
	RetrieveOperatingSystemInformation() (OperatingSystemInformation, error)
}

type SoftwareOperatingSystem struct{}

var _ OperatingSystem = (*SoftwareOperatingSystem)(nil)
