//go:build linux

package network

import (
	"fmt"
	"log/slog"
	"net"
)

type NetworkAdapterInformation struct {
	Description string `json:"Description"`
	MACAddress  string `json:"MACAddress"`
	IPAddress   string `json:"IPAddress"`
	DHCPEnabled string `json:"DHCPEnabled"`
}

func RetrieveActiveNetworkAdapters() ([]NetworkAdapterInformation, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		slog.Error("Failed to get network interfaces", slog.String("error_details", err.Error()))
		return nil, fmt.Errorf("failed to get interfaces: %w", err)
	}

	var adapters []NetworkAdapterInformation
	for _, iface := range interfaces {
		addrs, _ := iface.Addrs()
		ipStr := ""
		if len(addrs) > 0 {
			ipStr = addrs[0].String()
		}
		adapters = append(adapters, NetworkAdapterInformation{
			Description: iface.Name,
			MACAddress:  iface.HardwareAddr.String(),
			IPAddress:   ipStr,
			DHCPEnabled: "Unknown",
		})
	}
	return adapters, nil
}