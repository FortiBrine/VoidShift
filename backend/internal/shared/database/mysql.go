package database

import (
	"log/slog"

	"github.com/FortiBrine/VoidShift/internal/shared/config"
	"github.com/FortiBrine/VoidShift/internal/shared/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMysqlDatabase(config config.Config, l *slog.Logger) (*gorm.DB, error) {
	db, err := gorm.Open(
		mysql.Open(config.MysqlDsn), &gorm.Config{
			Logger: logger.NewGormLogger(l, config.Environment),
		},
	)

	if err != nil {
		return nil, err
	}

	return db, nil
}
