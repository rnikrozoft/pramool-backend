package migrations

import (
	"context"
	"embed"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

//go:embed *.sql
var sqlMigrations embed.FS

func GetMigrator(ctx context.Context, db *bun.DB) (*migrate.Migrator, error) {
	migrations := migrate.NewMigrations()
	if err := migrations.Discover(sqlMigrations); err != nil {
		return nil, err
	}
	return migrate.NewMigrator(db, migrations), nil
}

func Migrate(ctx context.Context, db *bun.DB) (*migrate.MigrationGroup, error) {
	migrator, err := GetMigrator(ctx, db)
	if err != nil {
		return nil, err
	}
	if err := migrator.Init(ctx); err != nil {
		return nil, err
	}
	group, err := migrator.Migrate(ctx)
	if err != nil {
		return nil, err
	}
	return group, nil
}

func Rollback(ctx context.Context, db *bun.DB) (*migrate.MigrationGroup, error) {
	migrator, err := GetMigrator(ctx, db)
	if err != nil {
		return nil, err
	}
	if err := migrator.Init(ctx); err != nil {
		return nil, err
	}
	migration, err := migrator.Rollback(ctx)
	if err != nil {
		return nil, err
	}
	return migration, nil
}
