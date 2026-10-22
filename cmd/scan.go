/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/risicoding/nethawk/internal/scanner"
	"github.com/spf13/cobra"
)

var jsonOutput bool

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Discover Nethawk servers on your local network",
	Long: `Scan the local network for reachable Nethawk servers.

Nethawk discovers advertised services using multicast DNS (mDNS) and
queries discovered instances to retrieve their server status, IP addresses,
and device configuration.

By default, results are printed in a human-readable format. Use --json
to produce indented JSON suitable for scripts, other programs, and
automated processing.

Discovery depends on network connectivity and multicast DNS availability.
Firewalls, wireless client isolation, VLAN boundaries, and other network
restrictions may prevent some devices from being discovered or reached.

This command discovers Nethawk instances; it does not scan every host
or port on the network.`,
	Example: `  # Discover servers and display readable results
  nethawk scan

  # Output results as JSON
  nethawk scan --json

  # Use the JSON output in another command
  nethawk scan --json > servers.json`,
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := scanner.Scan()
		if err != nil {
			return fmt.Errorf("scan network: %w", err)
		}

		if jsonOutput {
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "\t")
			if err := encoder.Encode(resp); err != nil {
				return fmt.Errorf("encode scan results as JSON: %w", err)
			}
			return nil
		}

		if len(resp) == 0 {
			fmt.Println("No Nethawk servers found.")
			return nil
		}

		for i, server := range resp {
			if i > 0 {
				fmt.Println()
			}

			fmt.Printf(
				"Server: %s\nStatus: %s\nIPs: %s\nDevice:\n\tName: %s\n\tUname: %s\n\tAlias: %s\n",
				server.Name,
				server.Status,
				strings.Join(server.Ips, " "),
				server.Config.Name,
				server.Config.Uname,
				server.Config.Alias,
			)
		}

		return nil
	},
}

func init() {
	scanCmd.Flags().BoolVarP(
		&jsonOutput,
		"json",
		"j",
		false,
		"Output scan results as indented JSON",
	)

	rootCmd.AddCommand(scanCmd)
}
