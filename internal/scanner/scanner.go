package scanner

import (
	"sync"

	"github.com/risicoding/nethawk/internal/cache"
	"github.com/risicoding/nethawk/internal/client"
	"github.com/risicoding/nethawk/internal/mdns"
	"github.com/risicoding/nethawk/internal/protocol"
)

type InstanceWithServers struct {
	Instance mdns.MDNSInstance
	Servers  []protocol.ServerWithLatency
}

func Scan() ([]protocol.ServerWithLatency, error) {
	instances, err := mdns.Scan()
	if err != nil {
		return nil, err
	}

	results := make([]InstanceWithServers, len(instances))

	var wg sync.WaitGroup

	for i, instance := range instances {
		wg.Add(1)

		go func(i int, ins mdns.MDNSInstance) {
			defer wg.Done()

			servers := make([]protocol.ServerWithLatency, len(ins.Ipv4))
			success := make([]bool, len(ins.Ipv4))

			var ipwg sync.WaitGroup

			for j, ipaddr := range ins.Ipv4 {
				ipwg.Add(1)

				go func(j int, ip string) {
					defer ipwg.Done()

					resp, err := client.GETROOT(ip)
					if err != nil {
						return
					}

					servers[j] = resp
					success[j] = true
				}(j, ipaddr.String())
			}

			ipwg.Wait()

			instanceServers := make(
				[]protocol.ServerWithLatency,
				0,
				len(servers),
			)

			for j, server := range servers {
				if success[j] {
					instanceServers = append(instanceServers, server)
				}
			}

			results[i] = InstanceWithServers{
				Instance: ins,
				Servers:  instanceServers,
			}
		}(i, instance)
	}

	wg.Wait()

	discovered := make([]protocol.ServerWithLatency, 0, len(results))

	for _, result := range results {
		if len(result.Servers) == 0 {
			continue
		}

		best := result.Servers[0]

		for _, server := range result.Servers[1:] {
			if server.Latency < best.Latency {
				best = server
			}
		}

		discovered = append(discovered, best)
	}

	return discovered, nil
}

var serverCache = cache.New(Scan)
var once sync.Once

func ScanCached() ([]protocol.ServerWithLatency, error) {
	once.Do(func() {
		serverCache.Refresh()
	})

	return serverCache.Get()
}
