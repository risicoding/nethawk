package ip

import (
	"errors"
	"fmt"
	"net"
	"sync"

	"github.com/risicoding/nethawk/internal/cache"
	"github.com/risicoding/nethawk/internal/client"
	internalNmap "github.com/risicoding/nethawk/internal/nmap"
	"github.com/risicoding/nethawk/internal/protocol"
)

type InterfaceWithIps struct {
	Interface Interface
	Ips       []string
}

type InterfaceWithServers struct {
	Interface Interface
	Servers   []protocol.Server
}

func ScanInterfaceWithClient(inf Interface) (InterfaceWithServers, error) {
	ips, err := GetIPRange(inf.Ipv4, inf.Mask)
	if err != nil {
		return InterfaceWithServers{}, err
	}

	var servers []protocol.Server

	ch := make(chan (protocol.Server))

	var wg sync.WaitGroup
	for _, ip := range ips {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			resp, err := client.Query(ip)
			if err != nil {
				return
			}

			ch <- resp
		}(ip)
	}

	wg.Wait()
	defer close(ch)

	for v := range ch {
		servers = append(servers, v)
	}

	return InterfaceWithServers{Interface: inf, Servers: servers}, nil
}

func ScanInterface(inf Interface) (InterfaceWithIps, error) {
	ip := inf.Ipv4
	mask := inf.Mask

	parsedIP := net.ParseIP(ip).To4()
	if parsedIP == nil {
		return InterfaceWithIps{}, fmt.Errorf(
			"invalid IPv4 address: %q", ip,
		)
	}

	parsedMask := net.ParseIP(mask).To4()
	if parsedMask == nil {
		return InterfaceWithIps{}, fmt.Errorf(
			"invalid IPv4 subnet mask: %q", mask,
		)
	}

	maskIP := net.IPMask(parsedMask)
	ones, bits := maskIP.Size()
	if bits != 32 {
		return InterfaceWithIps{}, fmt.Errorf(
			"invalid subnet mask: %q", mask,
		)
	}

	network := &net.IPNet{
		IP:   parsedIP.Mask(maskIP),
		Mask: maskIP,
	}

	if ones == 0 {
		return InterfaceWithIps{}, fmt.Errorf(
			"refusing to scan an entire IPv4 address space",
		)
	}

	ips, err := internalNmap.NmapScan(network.String())
	if err != nil {
		return InterfaceWithIps{}, err
	}

	return InterfaceWithIps{
		Interface: inf,
		Ips:       ips,
	}, nil
}

var (
	interfaceCacheStore = make(map[string]*cache.Cache[InterfaceWithServers])
	interfaceCacheMu    sync.Mutex
)

func ScanInterfaceCached(inf Interface) (InterfaceWithServers, error) {
	interfaceCacheMu.Lock()
	defer interfaceCacheMu.Unlock()

	ch, ok := interfaceCacheStore[inf.Name]
	if !ok {
		ch = cache.New(func() (InterfaceWithServers, error) {
			return ScanInterfaceWithClient(inf)
		})

		interfaceCacheStore[inf.Name] = ch
	}

	return ch.Get()
}

func implScanAllInterfaces(
	scanInterfaceFunc func(Interface) (InterfaceWithServers, error),
) ([]InterfaceWithServers, error) {
	infsData, err := GetAllInterfacesWithData()
	if err != nil {
		return nil, err
	}

	chanIps := make(chan InterfaceWithServers, len(infsData))
	chanErrors := make(chan error, len(infsData))

	var wg sync.WaitGroup

	for _, inf := range infsData {
		wg.Add(1)

		go func(inf Interface) {
			defer wg.Done()

			data, err := scanInterfaceFunc(inf)
			if err != nil {
				chanErrors <- fmt.Errorf(
					"scan interface %s (%s): %w",
					inf.Name, inf.Ipv4, err,
				)
				return
			}

			chanIps <- data
		}(inf)
	}

	wg.Wait()
	close(chanIps)
	close(chanErrors)

	var results []InterfaceWithServers
	for inf := range chanIps {
		results = append(results, inf)
	}

	var errs []error
	for err := range chanErrors {
		errs = append(errs, err)
	}

	return results, errors.Join(errs...)
}

func ScanAllInterfaces() ([]InterfaceWithServers, error) {
	return implScanAllInterfaces(ScanInterfaceWithClient)
}

func ScanAllInterfacesCached() ([]InterfaceWithServers, error) {
	return implScanAllInterfaces(ScanInterfaceCached)
}
