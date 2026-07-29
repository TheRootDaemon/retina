package main

import (
	"github.com/therootdaemon/hop/internal/config"
	"github.com/therootdaemon/hop/logger"
)

func main() {
	if err := config.Initialize(); err != nil {
		logger.Exit("failed to load config: %v", err)
	}

	logger.SetDefaultFromConfig(config.Logger())
	defer func() {
		_ = logger.Close()
	}()
}
