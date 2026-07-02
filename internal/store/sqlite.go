package store

import (
	"database/sql"

	"github.com/FortiBrine/VoidShift/internal/config"
	_ "modernc.org/sqlite"
)

func NewSqliteDatabase(cfg config.Config) (*sql.DB, error) {
	db, err := sql.Open("sqlite", cfg.SqliteDatabasePath)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}
