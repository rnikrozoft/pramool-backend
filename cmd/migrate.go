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

// migrateCmd applies pending SQL migrations embedded in the binary (*.up.sql).
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Apply pending database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "migrate target: %s\n", databaseTargetLine())
		fmt.Fprintln(cmd.OutOrStdout(), "         (bun tracks applied files in table bun_migrations)")
		group, err := migrations.Migrate(context.Background(), conn)
		if err != nil {
			panic(err)
		}
		if group == nil || group.IsZero() {
			fmt.Fprintln(cmd.OutOrStdout(), "migrate: already up to date (no pending migrations)")
			return
		}
		fmt.Fprintf(cmd.OutOrStdout(), "migrate OK: %s\n", group.String())
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// migrateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// migrateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
