package applications

type ApplicationInformation struct {
	DisplayName    string `json:"DisplayName"`
	DisplayVersion string `json:"DisplayVersion"`
	Publisher      string `json:"Publisher"`
	InstallDate    string `json:"InstallDate"`
}

type Applications interface {
	RetrieveInstalledApplications() ([]ApplicationInformation, error)
}

type SoftwareApplications struct{}

var _ Applications = (*SoftwareApplications)(nil)