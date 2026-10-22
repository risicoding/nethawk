/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
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
	deviceUname string
	deviceAlias string
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the nethawk server",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer cancel()

		logger.Info("after mdns")
		deviceUname := strings.Join(strings.Split(deviceName, " "), "_")
		device := protocol.DeviceConfig{
			Name:  deviceName,
			Uname: deviceUname,
			Alias: deviceAlias,
		}

		infs, err := ip.GetInterfacesWithIps()
		if err != nil {
			panic(err)
		}

		go mdns.Transmit(ctx, deviceName)
		go echo.Start(ctx, device, infs, scanner.ScanCached)

		<-ctx.Done()
		logger.Info("Stopping nethawk")
	},
}

func init() {
	startCmd.Flags().StringVarP(
		&deviceName,
		"name",
		"n",
		"",
		"Device name",
	)

	startCmd.Flags().StringVarP(
		&deviceAlias,
		"alias",
		"a",
		"",
		"Device alias",
	)

	_ = startCmd.MarkFlagRequired("name")
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
