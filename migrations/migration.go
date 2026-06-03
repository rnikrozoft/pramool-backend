package migrations

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strings"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

const (
	DBCore    = "core"
	DBWallet  = "wallet"
	DBAuction = "auction"
	DBAdmin   = "admin"
	DBAll     = "all"
)

//go:embed core/*.sql
var coreSQL embed.FS

//go:embed wallet/*.sql
var walletSQL embed.FS

//go:embed auction/*.sql
var auctionSQL embed.FS

//go:embed admin/*.sql
var adminSQL embed.FS

var allMigrationFS = []fs.FS{coreSQL, walletSQL, auctionSQL, adminSQL}

func discoverOn(fsys fs.FS) (*migrate.Migrations, error) {
	m := migrate.NewMigrations()
	if err := m.Discover(fsys); err != nil {
		return nil, err
	}
	return m, nil
}

func combineSorted(fsys []fs.FS) (*migrate.Migrations, error) {
	var all migrate.MigrationSlice
	for _, partFS := range fsys {
		part, err := discoverOn(partFS)
		if err != nil {
			return nil, err
		}
		all = append(all, part.Sorted()...)
	}
	sort.Slice(all, func(i, j int) bool {
		return all[i].Name < all[j].Name
	})
	combined := migrate.NewMigrations()
	for _, mig := range all {
		combined.Add(mig)
	}
	return combined, nil
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
	case DBAdmin:
		return discoverOn(adminSQL)
	case DBAll:
		return combineSorted(allMigrationFS)
	default:
		return nil, fmt.Errorf("unknown migration db %q (use core, wallet, auction, admin, or all)", target)
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
