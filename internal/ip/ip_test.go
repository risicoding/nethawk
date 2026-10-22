package ip

import (
	"fmt"
	"strings"
	"testing"
)

func TestInterfaces(t *testing.T) {
	infs, err := GetInterfacesWithIps()
	if err != nil {
		panic(err)
	}
	for _, inf := range infs {
		var ips []string
		for _, ip := range inf.IPs {
			ips = append(ips, ip.String())
		}

		fmt.Printf("Interface\n\tName: %s\n\tIps: %s\n", inf.Name, strings.Join(ips, " | "))
	}
}
