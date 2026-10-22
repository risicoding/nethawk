/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

	jsonlib "encoding/json"

	"github.com/risicoding/nethawk/internal/scanner"
	"github.com/spf13/cobra"
)

var json = false

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan the network for nethawk servers",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := scanner.Scan()
		if err != nil {
			panic(err)
		}

		if json {
			v, err := jsonlib.MarshalIndent(resp, "", "\t")
			if err != nil {
				panic(err)
			}
			fmt.Println(string(v))
			return
		}
		for _, server := range resp {
			fmt.Printf("Server: %s\nStatus: %s\nIps: %s\nDevice:\n\tName: %s\n\tUname: %s\n\tAlias: %s\n", server.Name, server.Status, strings.Join(server.Ips, " "), server.Config.Name, server.Config.Uname, server.Config.Alias)
		}
		// devices, err := mdns.Scan()
		// if err != nil {
		// 	panic(err)
		// }
		// for _, device := range devices {
		// 	var ips []string
		// 	for _, ip := range device.Ipv4 {
		// 		ips = append(ips, ip.String())
		// 	}
		//
		// 	fmt.Printf("NAME: %s\nINSTANCE: %s\nIPV4: %s\n", device.Hostname, device.Instance, strings.Join(ips, " "))
		// }
	},
}

func init() {
	scanCmd.Flags().BoolVarP(&json, "json", "j", false, "Output results in json")

	rootCmd.AddCommand(scanCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
