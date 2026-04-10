package database

import (
	"log/slog"

	"github.com/FortiBrine/VoidShift/internal/shared/config"
	"github.com/FortiBrine/VoidShift/internal/shared/logger"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func NewSqliteDatabase(config config.Config, l *slog.Logger) (*gorm.DB, error) {
	db, err := gorm.Open(
		sqlite.Open(config.SqliteDatabasePath), &gorm.Config{
			Logger: logger.NewGormLogger(l, config.Environment),
		},
	)

	if err != nil {
		return nil, err
	}

	return db, nil
}
