package migrations

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"strings"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

const (
	DBCore    = "core"
	DBWallet  = "wallet"
	DBAuction = "auction"
	DBAll     = "all"
)

//go:embed core/*.sql
var coreSQL embed.FS

//go:embed wallet/*.sql
var walletSQL embed.FS

//go:embed auction/*.sql
var auctionSQL embed.FS

func discoverOn(fsys fs.FS) (*migrate.Migrations, error) {
	m := migrate.NewMigrations()
	if err := m.Discover(fsys); err != nil {
		return nil, err
	}
	return m, nil
}

func migrationsFor(target string) (*migrate.Migrations, error) {
	target = strings.TrimSpace(strings.ToLower(target))
	if target == "" {
		target = DBAll
	}

	switch target {
	case DBCore:
		return discoverOn(coreSQL)
	case DBWallet:
		return discoverOn(walletSQL)
	case DBAuction:
		return discoverOn(auctionSQL)
	case DBAll:
		combined := migrate.NewMigrations()
		for _, fsys := range []fs.FS{coreSQL, walletSQL, auctionSQL} {
			part, err := discoverOn(fsys)
			if err != nil {
				return nil, err
			}
			for _, mig := range part.Sorted() {
				combined.Add(mig)
			}
		}
		return combined, nil
	default:
		return nil, fmt.Errorf("unknown migration db %q (use core, wallet, auction, or all)", target)
	}
}

func GetMigrator(ctx context.Context, db *bun.DB, target string) (*migrate.Migrator, error) {
	m, err := migrationsFor(target)
	if err != nil {
		return nil, err
	}
	return migrate.NewMigrator(db, m), nil
}

func Migrate(ctx context.Context, db *bun.DB, target string) (*migrate.MigrationGroup, error) {
	migrator, err := GetMigrator(ctx, db, target)
	if err != nil {
		return nil, err
	}
	if err := migrator.Init(ctx); err != nil {
		return nil, err
	}
	return migrator.Migrate(ctx)
}

func Rollback(ctx context.Context, db *bun.DB, target string) (*migrate.MigrationGroup, error) {
	migrator, err := GetMigrator(ctx, db, target)
	if err != nil {
		return nil, err
	}
	if err := migrator.Init(ctx); err != nil {
		return nil, err
	}
	return migrator.Rollback(ctx)
}
