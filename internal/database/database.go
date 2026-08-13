package database

import (
	"github.com/therootdaemon/hop/internal/model"
	"github.com/therootdaemon/hop/logger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// migrate applies the database schema migrations
// for the application's models.
var migrate = func(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.Host{},
	)
}

// Connect opens a database connection
// and runs the database migrations.
//
// If the database cannot be opened or the migrations fail,
// Connect returns the corresponding error and no database handle.
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
