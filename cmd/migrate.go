/*
Copyright © 2025 rnikrozoft rnikrozoft.dev@gmail.com
*/
package cmd

import (
	"context"
	"fmt"

	"github.com/rnikrozoft/pramool-core/migrations"
	"github.com/spf13/cobra"
)

var migrateDB string

// migrateCmd applies pending SQL migrations for core, wallet, and/or auction databases.
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Apply pending database migrations",
	Long: `Apply SQL migrations from pramool-core/migrations/{core,wallet,auction}/.
All targets use the same PostgreSQL database (DATABASE_NAME / DATABASE_DSN).

Examples:
  go run . migrate --db core
  go run . migrate --db wallet
  go run . migrate --db auction
  go run . migrate --db all    # core + wallet + auction SQL on one database`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		target := migrateDB
		if target == migrations.DBAll {
			if err := runMigrate(ctx, cmd, migrations.DBAll); err != nil {
				panic(err)
			}
			return
		}
		if err := runMigrate(ctx, cmd, target); err != nil {
			panic(err)
		}
	},
}

func runMigrate(ctx context.Context, cmd *cobra.Command, target string) error {
	db, label, err := openMigrateDB(target)
	if err != nil {
		return err
	}
	defer db.Close()

	fmt.Fprintf(cmd.OutOrStdout(), "migrate: %s (scope=%s)\n", label, target)
	fmt.Fprintln(cmd.OutOrStdout(), "         (bun tracks applied files in table bun_migrations)")
	group, err := migrations.Migrate(ctx, db, target)
	if err != nil {
		return err
	}
	if group == nil || group.IsZero() {
		fmt.Fprintf(cmd.OutOrStdout(), "migrate [%s]: already up to date\n", target)
		return nil
	}
	fmt.Fprintf(cmd.OutOrStdout(), "migrate [%s] OK: %s\n", target, group.String())
	return nil
}

func init() {
	migrateCmd.Flags().StringVar(&migrateDB, "db", migrations.DBAll, "database: core, wallet, auction, or all")
	rootCmd.AddCommand(migrateCmd)
}
