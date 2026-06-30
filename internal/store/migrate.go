package store

import (
	"context"
	"database/sql"

	"github.com/FortiBrine/VoidShift/db/migrations"
	"github.com/pressly/goose/v3"
)

func Migrate(ctx context.Context, db *sql.DB) error {
	provider, err := goose.NewProvider(goose.DialectSQLite3, db, migrations.FS)
	if err != nil {
		return err
	}
	_, err = provider.Up(ctx)
	return err
}
