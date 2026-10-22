package main

import (
	"github.com/risicoding/nethawk/cmd"
	"github.com/risicoding/nethawk/internal/logger"
)

func main() {
	cleanup := logger.Init()
	defer cleanup()
	cmd.Execute()
}
