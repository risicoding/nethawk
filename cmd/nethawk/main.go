package main

import (
	"fmt"

	//server "github.com/risicoding/nethawk/cmd"
	"github.com/risicoding/nethawk/internal/ip"
	"github.com/risicoding/nethawk/internal/server"
)

func Scanner() {
	data, err := ip.ScanAllInterfacesCached()

	if err != nil {
		panic(err)
	}

	for _, inf := range data {
		fmt.Printf("Interface:\n\n")
		fmt.Printf("\tName: %s\n", inf.Interface.Name)
		fmt.Printf("\tIpv4: %s\n", inf.Interface.Ipv4)
		fmt.Printf("\n")

		fmt.Printf("Ip addresses:\n\n")
		for _, server := range inf.Servers {
			fmt.Printf("\t%+v\n", server)
		}
		fmt.Printf("\n")
	}
}

func main() {
	go server.Start()

	for i := range 20 {
		fmt.Printf("iteration %d\n", i)

		v, err := ip.ScanAllInterfacesCached()
		if err != nil {
			panic(err)
		}

		for _, inf := range v {
			fmt.Printf("Interface %s %s\n", inf.Interface.Name, inf.Interface.Ipv4)
			for _, server := range inf.Servers {

				fmt.Printf("Server %s\n", server.Name)
				fmt.Printf("Device %s\n", server.DeviceName)
				fmt.Printf("DeviceUname %s\n", server.DeviceUname)
				fmt.Printf("DeviceAlias %s\n", server.DeviceAlias)

			}
		}
		fmt.Printf("%+v\n", v)
	}
}
