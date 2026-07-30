package network

import (
	"net"
)

type NetworkInformation struct {
	Name         string   `json:"Name"`
	HardwareAddr string   `json:"HardwareAddr"`
	IPAddresses  []string `json:"IPAddresses"`
	Flags        string   `json:"Flags"`
}

func RetrieveNetworkInfo() ([]NetworkInformation, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var netInfos []NetworkInformation
	for _, iface := range interfaces {
		var ips []string
		addrs, err := iface.Addrs()
		if err == nil {
			for _, addr := range addrs {
				ips = append(ips, addr.String())
			}
		}

		netInfos = append(netInfos, NetworkInformation{
			Name:         iface.Name,
			HardwareAddr: iface.HardwareAddr.String(),
			IPAddresses:  ips,
			Flags:        iface.Flags.String(),
		})
	}

	return netInfos, nil
}