package services

type ServiceInformation struct {
	Name        string `json:"Name"`
	DisplayName string `json:"DisplayName"`
	State       string `json:"State"`
	StartMode   string `json:"StartMode"`
	StartName   string `json:"StartName"`
}

type Services interface {
	RetrieveSystemServices() ([]ServiceInformation, error)
}

type SoftwareServices struct{}

var _ Services = (*SoftwareServices)(nil)