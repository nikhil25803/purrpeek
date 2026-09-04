package network

import (
	"errors"
	"fmt"
	"net"
	"os"
	"slices"
	"strings"
)

const (
	defaultIPv4Destination = "192.0.2.1"
	defaultIPv6Destination = "2001:db8::1"
)

type NetworkInterface struct {
	Name       string   `json:"name"`
	Addresses  []string `json:"addresses"`
	MACAddress string   `json:"macAddress,omitempty"`
	MTU        int      `json:"mtu"`
}

type NetworkInfo struct {
	Hostname         string             `json:"hostname,omitempty"`
	PrimaryInterface string             `json:"primaryInterface,omitempty"`
	LocalIPv4        string             `json:"localIPv4,omitempty"`
	LocalIPv6        string             `json:"localIPv6,omitempty"`
	MACAddress       string             `json:"macAddress,omitempty"`
	Interfaces       []NetworkInterface `json:"interfaces"`
}

func GetNetworkInformation() (*NetworkInfo, error) {
	info := &NetworkInfo{Interfaces: []NetworkInterface{}}
	var errs []error
	if hostname, err := os.Hostname(); err != nil {
		errs = append(errs, fmt.Errorf("hostname: %w", err))
	} else {
		info.Hostname = hostname
	}

	interfaces, err := net.Interfaces()
	if err != nil {
		errs = append(errs, err)
		return info, errors.Join(errs...)
	}

	skipped := 0
	for _, networkInterface := range interfaces {
		if networkInterface.Flags&net.FlagUp == 0 || networkInterface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addresses, err := networkInterface.Addrs()
		if err != nil {
			skipped++
			continue
		}
		usableAddresses := make([]string, 0, len(addresses))
		for _, address := range addresses {
			ip, _, err := net.ParseCIDR(address.String())
			if err == nil && ip.IsGlobalUnicast() && !ip.IsLinkLocalUnicast() {
				usableAddresses = append(usableAddresses, address.String())
			}
		}
		if len(usableAddresses) == 0 {
			continue
		}
		slices.Sort(usableAddresses)
		info.Interfaces = append(info.Interfaces, NetworkInterface{
			Name:       networkInterface.Name,
			Addresses:  usableAddresses,
			MACAddress: networkInterface.HardwareAddr.String(),
			MTU:        networkInterface.MTU,
		})
	}
	if skipped > 0 {
		errs = append(errs, fmt.Errorf("%d interface(s) unavailable", skipped))
	}
	slices.SortFunc(info.Interfaces, func(a, b NetworkInterface) int {
		return strings.Compare(a.Name, b.Name)
	})

	if primary := findPrimaryInterface(info.Interfaces); primary != nil {
		info.PrimaryInterface = primary.Name
		info.MACAddress = primary.MACAddress
		for _, address := range primary.Addresses {
			ip, _, err := net.ParseCIDR(address)
			if err != nil {
				continue
			}
			if ip.To4() != nil && info.LocalIPv4 == "" {
				info.LocalIPv4 = ip.String()
			} else if ip.To4() == nil && info.LocalIPv6 == "" {
				info.LocalIPv6 = ip.String()
			}
		}
	}
	return info, errors.Join(errs...)
}

func findPrimaryInterface(interfaces []NetworkInterface) *NetworkInterface {
	for _, routeIP := range []net.IP{
		localRouteIP("udp4", defaultIPv4Destination),
		localRouteIP("udp6", defaultIPv6Destination),
	} {
		for index := range interfaces {
			for _, address := range interfaces[index].Addresses {
				ip, _, err := net.ParseCIDR(address)
				if err == nil && routeIP != nil && routeIP.Equal(ip) {
					return &interfaces[index]
				}
			}
		}
	}
	if len(interfaces) > 0 {
		return &interfaces[0]
	}
	return nil
}

func localRouteIP(network, destination string) net.IP {
	connection, err := net.DialUDP(network, nil, &net.UDPAddr{IP: net.ParseIP(destination), Port: 9})
	if err != nil {
		return nil
	}
	defer connection.Close()
	return connection.LocalAddr().(*net.UDPAddr).IP
}
