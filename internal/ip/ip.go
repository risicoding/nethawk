package ip

import (
	"errors"
	"fmt"
	"net"
	"net/netip"
)

type Interface struct {
	Name string
	Ipv4 string
	Mask string
}

var (
	ErrorNoIpv4Address = errors.New("interface has no ipv4 addresses")
	ErrorVirtInterface = errors.New("interface is not a physical interface")
)

func GetInterfaces() ([]net.Interface, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to find interfaces: %w", err)
	}

	return interfaces, nil
}

func GetInterfaceData(name string) (Interface, error) {
	if name == "lo" || name == "docker0" {
		return Interface{}, ErrorVirtInterface
	}

	inf, err := net.InterfaceByName(name)
	if err != nil {
		return Interface{}, err
	}

	addrs, err := inf.Addrs()
	if err != nil {
		return Interface{}, err
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok {
			if ipv4 := ipNet.IP.To4(); ipv4 != nil {
				mask := net.IP(ipNet.Mask)

				return Interface{
					Name: name,
					Ipv4: ipv4.String(),
					Mask: mask.String(),
				}, nil
			}
		}
	}

	return Interface{}, ErrorNoIpv4Address
}

func GetAllInterfacesWithData() ([]Interface, error) {
	infs, err := GetInterfaces()
	if err != nil {
		return nil, fmt.Errorf("get interfaces: %w", err)
	}

	var infsData []Interface

	for _, inf := range infs {
		data, err := GetInterfaceData(inf.Name)
		if err != nil {
			if errors.Is(err, ErrorNoIpv4Address) ||
				errors.Is(err, ErrorVirtInterface) {
				continue
			}

			return nil, fmt.Errorf(
				"get interface data for %s: %w",
				inf.Name, err,
			)
		}

		infsData = append(infsData, data)
	}

	return infsData, nil
}

func GetIPRange(ip, mask string) ([]string, error) {
	parsedIp := net.ParseIP(ip)
	parsedMask := net.ParseIP(mask)

	if parsedIp == nil || parsedMask == nil {
		return nil, fmt.Errorf("invalid ip or subnet mask")
	}

	maskBytes := net.IPMask(parsedMask.To4())
	ones, bits := maskBytes.Size()
	if bits != 32 {
		return nil, fmt.Errorf("invalid IPv4 subnet mask: %q", mask)
	}

	addr, err := netip.ParseAddr(ip)
	if err != nil {
		return nil, err
	}

	prefix := netip.PrefixFrom(addr, ones)
	currentIp := prefix.Masked().Addr()

	var ipRange []string

	for prefix.Contains(currentIp) {
		ipRange = append(ipRange, currentIp.String())
		currentIp = currentIp.Next()
	}

	return ipRange, nil
}
