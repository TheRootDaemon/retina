package database

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/therootdaemon/hop/internal/model"
	"gorm.io/gorm"
)

func TestConnect(t *testing.T) {
	oldMigrate := migrate
	t.Cleanup(
		func() {
			migrate = oldMigrate
		},
	)

	tests := []struct {
		name    string
		path    string
		migrate func(db *gorm.DB) error
		wantErr bool
	}{
		{
			name:    "database open failure",
			path:    t.TempDir(),
			migrate: func(db *gorm.DB) error { return nil },
			wantErr: true,
		},
		{
			name:    "in memory sqlite",
			path:    ":memory:",
			migrate: func(db *gorm.DB) error { return nil },
		},
		{
			name:    "migration failure",
			path:    ":memory:",
			migrate: func(db *gorm.DB) error { return errors.New("migrate failed") },
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			migrate = tt.migrate

			db, err := Connect(tt.path)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, db)
			} else {
				require.NoError(t, err)
				require.NotNil(t, db)

				sqlDB, err := db.DB()
				require.NoError(t, err)
				assert.NoError(t, sqlDB.Ping())
				assert.NoError(t, sqlDB.Close())
			}
		})
	}
}

func TestConnect_UsesDefaultMigrate(t *testing.T) {
	db, err := Connect(":memory:")
	require.NoError(t, err)
	require.NotNil(t, db)

	assert.True(t, db.Migrator().HasTable(&model.Host{}))
}
