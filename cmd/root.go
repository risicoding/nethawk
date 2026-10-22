/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "nethawk",
	Short: "Discover and inspect Nethawk instances on your local network",
	Long: `Nethawk is a local network discovery and device inspection CLI.

It discovers running Nethawk instances on your LAN using multicast DNS
(mDNS) and queries discovered servers for their device information.

Use Nethawk to identify devices running the Nethawk service, inspect their
network addresses, and retrieve their advertised device configuration.

Commands:
  start   Start a Nethawk server on this device
  scan    Discover Nethawk servers on the local network

Nethawk is designed for local network discovery. Devices must be reachable
on the network, and multicast DNS traffic must be permitted by the network
for discovery to work reliably.`,
	Version: "v0.1.0",
	Example: `  # Display available commands and global options
  nethawk --help

  # Discover Nethawk servers on your LAN
  nethawk scan

  # Discover servers and output machine-readable JSON
  nethawk scan --json

  # Start a Nethawk server with a device name
  nethawk start --name fedora-pc

  # Start a server with a human-friendly alias
  nethawk start --name fedora-pc --alias "My Desktop"`,
	SilenceUsage: true,
}

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Configure persistent flags here when global configuration options
	// are introduced.
}
