package network

import (
	"net"
	"os"
)

const (
	defaultIPv4Destination = "192.0.2.1"
	defaultIPv6Destination = "2001:db8::1"
)

type NetworkInterface struct {
	Name       string
	Addresses  []string
	MACAddress string
	MTU        int
}

type NetworkInfo struct {
	Hostname         string
	PrimaryInterface string
	LocalIPv4        string
	LocalIPv6        string
	MACAddress       string
	Interfaces       []NetworkInterface
}

func GetNetworkInformation() *NetworkInfo {
	info := &NetworkInfo{Interfaces: []NetworkInterface{}}
	info.Hostname, _ = os.Hostname()

	interfaces, err := net.Interfaces()
	if err != nil {
		return info
	}

	for _, networkInterface := range interfaces {
		if networkInterface.Flags&net.FlagUp == 0 || networkInterface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addresses, err := networkInterface.Addrs()
		if err != nil {
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

		info.Interfaces = append(info.Interfaces, NetworkInterface{
			Name:       networkInterface.Name,
			Addresses:  usableAddresses,
			MACAddress: networkInterface.HardwareAddr.String(),
			MTU:        networkInterface.MTU,
		})
	}

	if primary := findPrimaryInterface(info.Interfaces); primary != nil {
		info.PrimaryInterface = primary.Name
		info.MACAddress = primary.MACAddress
		for _, address := range primary.Addresses {
			ip, _, _ := net.ParseCIDR(address)
			if ip.To4() != nil && info.LocalIPv4 == "" {
				info.LocalIPv4 = ip.String()
			} else if ip.To4() == nil && info.LocalIPv6 == "" {
				info.LocalIPv6 = ip.String()
			}
		}
	}

	return info
}

func findPrimaryInterface(interfaces []NetworkInterface) *NetworkInterface {
	for _, routeIP := range []net.IP{
		localRouteIP("udp4", defaultIPv4Destination),
		localRouteIP("udp6", defaultIPv6Destination),
	} {
		for index := range interfaces {
			for _, address := range interfaces[index].Addresses {
				ip, _, _ := net.ParseCIDR(address)
				if routeIP != nil && routeIP.Equal(ip) {
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
	connection, err := net.DialUDP(network, nil, &net.UDPAddr{
		IP:   net.ParseIP(destination),
		Port: 9,
	})
	if err != nil {
		return nil
	}
	defer connection.Close()

	return connection.LocalAddr().(*net.UDPAddr).IP
}
