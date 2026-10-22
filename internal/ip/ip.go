package ip

import (
	"fmt"
	"net"
	"strings"

	"github.com/wlynxg/anet"
)

func GetInterfaces() ([]net.Interface, error) {
	infs, err := anet.Interfaces()
	if err != nil {
		return nil, err
	}

	var interfaces []net.Interface
	for _, inf := range infs {
		if inf.Name == "docker0" || inf.Name == "lo" || strings.Contains(inf.Name, "ccmni") {
			continue
		}

		addrs, err := inf.Addrs()
		if err != nil {
			continue
		}

		exists := false
		for _, addr := range addrs {

			if ipNet, ok := addr.(*net.IPNet); ok {

				if ipv4 := ipNet.IP.To4(); ipv4 != nil {
					exists = true

				}
			}

		}
		if exists {
			interfaces = append(interfaces, inf)
		}
	}

	return interfaces, nil
}

type InterfaceWithIPs struct {
	net.Interface
	IPs []net.IP
}

func GetInterfacesWithIps() ([]InterfaceWithIPs, error) {
	infs, err := anet.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("error with anet interfaces %w", err)
	}

	var interfaces []InterfaceWithIPs
	for _, inf := range infs {
		if inf.Name == "docker0" || inf.Name == "lo" || strings.Contains(inf.Name, "ccmni") {
			continue
		}

		addrs, err := inf.Addrs()
		if err != nil {
			continue
		}

		exists := false
		var ips []net.IP
		for _, addr := range addrs {

			if ipNet, ok := addr.(*net.IPNet); ok {

				if ipv4 := ipNet.IP.To4(); ipv4 != nil {
					exists = true

					ips = append(ips, ipv4)
				}
			}

		}

		if exists {
			interfaces = append(interfaces, InterfaceWithIPs{Interface: inf, IPs: ips})
		}
	}

	return interfaces, nil
}
