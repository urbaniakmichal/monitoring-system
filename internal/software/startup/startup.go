package startup

type StartupCommandInformation struct {
	Name     string `json:"Name"`
	Command  string `json:"Command"`
	Location string `json:"Location"`
	User     string `json:"User"`
}

type Startup interface {
	RetrieveStartupCommands() ([]StartupCommandInformation, error)
}

type SoftwareStartup struct{}

var _ Startup = (*SoftwareStartup)(nil)