package nmap

import (
	"fmt"

	"github.com/Ullaakut/nmap"
)

func NmapScan(netStr string) ([]string, error) {

	scanner, err := nmap.NewScanner(
		nmap.WithTargets(netStr),
		nmap.WithPingScan(),
		nmap.WithTimingTemplate(nmap.TimingAggressive),
	)
	if err != nil {
		return nil, fmt.Errorf("create nmap scanner: %w", err)
	}

	result, _, err := scanner.Run()
	if err != nil {
		return nil, fmt.Errorf("run ping sweep: %w", err)
	}

	ips := make([]string, 0, len(result.Hosts))
	for _, host := range result.Hosts {
		if host.Status.State != "up" {
			continue
		}

		for _, addr := range host.Addresses {
			if addr.AddrType == "ipv4" {
				ips = append(ips, addr.Addr)
				break
			}
		}
	}

	return ips, nil
}
