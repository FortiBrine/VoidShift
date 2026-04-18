package database

import (
	"errors"
	"log/slog"

	"github.com/FortiBrine/VoidShift/internal/shared/config"
	"gorm.io/gorm"
)

func Open(cfg config.Config, l *slog.Logger) (*gorm.DB, error) {
	switch {
	case cfg.SqliteDatabasePath != "":
		return NewSqliteDatabase(cfg, l)
	case cfg.MysqlDsn != "":
		return NewMysqlDatabase(cfg, l)
	default:
		return nil, errors.New("no database configured")
	}
}
