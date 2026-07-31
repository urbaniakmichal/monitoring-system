package startup

type StartupCommandInformation struct {
	Name     string `json:"Name"`
	Command  string `json:"Command"`
	Location string `json:"Location"`
	User     string `json:"User"`
}