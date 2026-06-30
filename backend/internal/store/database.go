package store

import (
	"database/sql"
	"errors"

	"github.com/FortiBrine/VoidShift/internal/config"
)

func Open(cfg config.Config) (*sql.DB, error) {
	switch {
	case cfg.SqliteDatabasePath != "":
		return NewSqliteDatabase(cfg)
	default:
		return nil, errors.New("no store configured")
	}
}
