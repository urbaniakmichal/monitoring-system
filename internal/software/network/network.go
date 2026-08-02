package network

type NetworkAdapterInformation struct {
	Caption     string `json:"Caption"`
	Description string `json:"Description"`
	IPAddress   string `json:"IPAddress"`
	MACAddress  string `json:"MACAddress"`
	DHCPEnabled string `json:"DHCPEnabled"`
}

type Networks interface {
	RetrieveActiveNetworkAdapters() ([]NetworkAdapterInformation, error)
}

type SoftwareNetworks struct{}

var _ Networks = (*SoftwareNetworks)(nil)