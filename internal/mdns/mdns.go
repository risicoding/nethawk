package mdns

import (
	"context"
	"net"
	"strings"
	"time"

	"fmt"

	"github.com/grandcat/zeroconf"
	"github.com/risicoding/nethawk/internal/logger"
	"github.com/wlynxg/anet"
)

const PORT = 5354

type MDNSInstance struct {
	Hostname string
	Instance string
	Ipv4     []net.IP
	Port     int
}

func Transmit(ctx context.Context, deviceName string) error {
	infs, err := anet.Interfaces()
	if err != nil {
		return err
	}

	var interfaces []net.Interface

	for _, inf := range infs {
		if inf.Name == "lo" || inf.Name == "docker0" || strings.Contains(inf.Name, "ccmni") {
			continue
		}

		addrs, err := inf.Addrs()
		if err != nil {
			continue
		}

		hasIpv4 := false
		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok {
				if ipv4 := ipNet.IP.To4(); ipv4 != nil {

					hasIpv4 = true

				}
			}
		}

		if hasIpv4 {
			interfaces = append(interfaces, inf)
		}
	}

	server, err := zeroconf.Register(
		deviceName,
		"_nethawk._tcp",
		"local.",
		PORT,
		[]string{"version=1"},
		interfaces,
	)
	if err != nil {
		return fmt.Errorf("register mDNS service: %w", err)
	}
	defer server.Shutdown()

	logger.Infof("Advertising device: %s", deviceName)
	logger.Info("Service: _nethawk._tcp.local.")
	logger.Infof("Application port: %d", PORT)

	<-ctx.Done()
	logger.Warn("Stopping advertisement.")
	return nil
}

func Scan() ([]MDNSInstance, error) {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return nil, fmt.Errorf("create mDNS resolver: %w", err)
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	entries := make(chan *zeroconf.ServiceEntry)

	if err := resolver.Browse(
		ctx,
		"_nethawk._tcp",
		"local.",
		entries,
	); err != nil {
		return nil, fmt.Errorf("browse for services: %w", err)
	}

	found := make(map[string]struct{})
	var devices []MDNSInstance

	for {
		select {
		case <-ctx.Done():
			return devices, nil

		case entry := <-entries:
			if entry == nil {
				continue
			}

			// Keep only IPv4 addresses.
			ipv4 := make([]net.IP, 0, len(entry.AddrIPv4))
			for _, ip := range entry.AddrIPv4 {
				if addr := ip.To4(); addr != nil {
					ipv4 = append(ipv4, addr)
				}
			}

			// Ignore services without IPv4 addresses.
			if len(ipv4) == 0 {
				continue
			}

			// Distinguish services by instance, hostname, and port.
			key := fmt.Sprintf(
				"%s|%s|%d",
				entry.Instance,
				entry.HostName,
				entry.Port,
			)
			if _, exists := found[key]; exists {
				continue
			}
			found[key] = struct{}{}

			devices = append(devices, MDNSInstance{
				Instance: entry.Instance,
				Hostname: entry.HostName,
				Port:     entry.Port,
				Ipv4:     ipv4,
			})
		}
	}
}
