package database

import (
	"github.com/FortiBrine/VoidShift/internal/shared/config"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func NewSqliteDatabase(config config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(
		sqlite.Open(config.SqliteDatabasePath), &gorm.Config{},
	)

	if err != nil {
		return nil, err
	}

	return db, nil
}
