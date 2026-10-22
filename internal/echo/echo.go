package echo

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/risicoding/nethawk/internal/ip"
	"github.com/risicoding/nethawk/internal/logger"
	"github.com/risicoding/nethawk/internal/mdns"
	"github.com/risicoding/nethawk/internal/protocol"
)

var (
	device     protocol.DeviceConfig
	interfaces []ip.InterfaceWithIPs
	ips        []string
	scanner    func() ([]protocol.ServerWithLatency, error)
)

const PORT = 6363

func Initialize(deviceConfig protocol.DeviceConfig, infs []ip.InterfaceWithIPs, scannerFunc func() ([]protocol.ServerWithLatency, error)) {
	device = deviceConfig
	interfaces = infs
	scanner = scannerFunc

	var ipaddrs []string
	for _, inf := range infs {
		var infips []string
		for _, ip := range inf.IPs {

			infips = append(infips, ip.String())
		}

		for _, ip := range infips {
			ipaddrs = append(ipaddrs, ip)
		}
	}

}

func Start(ctx context.Context, deviceConfig protocol.DeviceConfig, infs []ip.InterfaceWithIPs, scannerFunc func() ([]protocol.ServerWithLatency, error)) {
	Initialize(
		deviceConfig,
		infs,
		scannerFunc,
	)

	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(LoggerMiddleware())

	e.GET("/", rootHandler)
	e.GET("/scan", scanHandler)
	e.GET("/interfaces", interfaceHandler)

	// 1. Create a v5 StartConfig structure
	sc := echo.StartConfig{
		Address:         fmt.Sprintf(":%d", PORT),
		GracefulTimeout: 5 * time.Second, // Max time allowed for mid-flight requests to finish
	}

	logger.Infof("Server starting on port :%d", PORT)

	// 2. Start the server using the configuration block.
	// This takes care of blocking and handling your `ctx` signal automatically!
	if err := sc.Start(ctx, e); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("Server encountered an error: ", err)
	}

	logger.Info("Server stopped")
}

func rootHandler(c *echo.Context) error {
	return c.JSON(http.StatusOK, protocol.Server{
		Name:   "nethawk",
		Status: "running",
		Config: device,
		Ips:    ips,
	})
}

func interfaceHandler(c *echo.Context) error {
	return c.JSON(http.StatusOK, interfaces)
}

func mdnsHandler(c *echo.Context) error {
	resp, err := mdns.Scan()
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to scan with mdns")
	}
	return c.JSON(http.StatusOK, resp)
}

func scanHandler(c *echo.Context) error {
	resp, err := scanner()
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to scan for nethawk instances")
	}

	return c.JSON(http.StatusOK, resp)
}
