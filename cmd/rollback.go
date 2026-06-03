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

// rollbackCmd rolls back the last applied SQL migration group for the selected database.
var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Roll back the last applied database migration group",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		db, label, err := openMigrateDB(rollbackDB)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		fmt.Fprintf(cmd.OutOrStdout(), "rollback target: %s (db=%s)\n", label, rollbackDB)
		fmt.Fprintln(cmd.OutOrStdout(), "            (bun tracks applied files in table bun_migrations)")
		group, err := migrations.Rollback(ctx, db, rollbackDB)
		if err != nil {
			panic(err)
		}
		if group == nil || group.IsZero() {
			fmt.Fprintf(cmd.OutOrStdout(), "rollback: nothing to roll back (no applied migrations)")
			return
		}
		fmt.Fprintf(cmd.OutOrStdout(), "rollback OK: %s\n", group.String())
	},
}

var rollbackDB string

func init() {
	rollbackCmd.Flags().StringVar(&rollbackDB, "db", migrations.DBCore, "database: core, wallet, auction, admin, or all")
	rootCmd.AddCommand(rollbackCmd)
}
