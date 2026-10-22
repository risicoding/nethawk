/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/risicoding/nethawk/internal/echo"
	"github.com/risicoding/nethawk/internal/ip"
	"github.com/risicoding/nethawk/internal/logger"
	"github.com/risicoding/nethawk/internal/mdns"
	"github.com/risicoding/nethawk/internal/protocol"
	"github.com/risicoding/nethawk/internal/scanner"
	"github.com/spf13/cobra"
)

var (
	deviceName  string
	deviceAlias string
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a Nethawk server on this device",
	Long: `Start the Nethawk service and advertise this device on the local network.

The server publishes its presence using multicast DNS (mDNS) and starts
the Nethawk discovery endpoint so other devices can identify this instance
and retrieve its status and device configuration.

The device name is required and identifies this instance in discovery
results. A machine-friendly identifier (uname) is derived from the device
name by replacing spaces with underscores.

An optional alias provides a human-friendly label for the device. For
example, a device named "fedora-pc" could use the alias "My Desktop".

The process runs in the foreground until it receives an interrupt or
termination signal, such as Ctrl+C or SIGTERM. On shutdown, the command
cancels its service context.

For discovery to work, the devices should share a network that permits
mDNS multicast traffic, and the server's listening endpoint must be
reachable through the local network.`,
	Example: `  # Start with a device name
  nethawk start --name fedora-pc

  # Set a friendly alias
  nethawk start --name fedora-pc --alias "My Desktop"

  # Use short flag names
  nethawk start -n laptop -a "Work Laptop"

  # Stop the foreground server
  # Press Ctrl+C`,
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := signal.NotifyContext(
			context.Background(),
			os.Interrupt,
			syscall.SIGTERM,
		)
		defer cancel()

		name := strings.TrimSpace(deviceName)
		if name == "" {
			return fmt.Errorf("device name cannot be empty")
		}

		uname := strings.Join(strings.Fields(name), "_")
		device := protocol.DeviceConfig{
			Name:  name,
			Uname: uname,
			Alias: strings.TrimSpace(deviceAlias),
		}

		infs, err := ip.GetInterfacesWithIps()
		if err != nil {
			return fmt.Errorf("get network interfaces: %w", err)
		}

		logger.Info("Starting Nethawk")

		go mdns.Transmit(ctx, name)
		go echo.Start(ctx, device, infs, scanner.ScanCached)

		<-ctx.Done()
		logger.Info("Stopping Nethawk")
		return nil
	},
}

func init() {
	startCmd.Flags().StringVarP(
		&deviceName,
		"name",
		"n",
		"",
		"Required device name advertised on the local network",
	)

	startCmd.Flags().StringVarP(
		&deviceAlias,
		"alias",
		"a",
		"",
		"Optional human-friendly alias for this device",
	)

	if err := startCmd.MarkFlagRequired("name"); err != nil {
		panic(err)
	}

	rootCmd.AddCommand(startCmd)
}
