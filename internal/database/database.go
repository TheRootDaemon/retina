package database

import (
	"github.com/therootdaemon/hop/internal/model"
	"github.com/therootdaemon/hop/logger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var migrate = func(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.Host{},
	)
}

func Connect(path string) (*gorm.DB, error) {
	logger.Debug("connecting to database at %s", path)

	db, err := gorm.Open(sqlite.Open(path))
	if err != nil {
		logger.Error("failed to open database: %v", err)
		return nil, err
	}
	logger.Debug("database connection established")

	if err := migrate(db); err != nil {
		logger.Error("database migration failed: %v", err)
		return nil, err
	}
	logger.Debug("database migration completed")

	return db, nil
}
