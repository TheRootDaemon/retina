package main

import (
	"github.com/therootdaemon/hop/internal/config"
	"github.com/therootdaemon/hop/internal/database"
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

	db, err := database.Connect(config.Database().Path)
	if err != nil {
		logger.Exit("failed to connect to database: %v", err)
	}
	_ = db
}
