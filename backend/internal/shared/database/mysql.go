package database

import (
	"github.com/FortiBrine/VoidShift/internal/shared/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMysqlDatabase(config config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(
		mysql.Open(config.MysqlDsn), &gorm.Config{},
	)

	if err != nil {
		return nil, err
	}

	return db, nil
}
