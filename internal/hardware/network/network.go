package network

type NetworkInformation struct {
	Name         string   `json:"Name"`
	HardwareAddr string   `json:"HardwareAddr"`
	IPAddresses  []string `json:"IPAddresses"`
	Flags        string   `json:"Flags"`
}

type Network interface {
	RetrieveNetworkInfo() ([]NetworkInformation, error)
}

type HardwareNetwork struct{}

var _ Network = (*HardwareNetwork)(nil)