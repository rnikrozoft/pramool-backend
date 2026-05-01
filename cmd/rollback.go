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

// rollbackCmd rolls back the last applied SQL migration group (runs paired *.down.sql).
var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Roll back the last applied database migration group",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "rollback target: %s\n", databaseTargetLine())
		fmt.Fprintln(cmd.OutOrStdout(), "            (bun tracks applied files in table bun_migrations)")
		group, err := migrations.Rollback(context.Background(), conn)
		if err != nil {
			panic(err)
		}
		if group == nil || group.IsZero() {
			fmt.Fprintln(cmd.OutOrStdout(), "rollback: nothing to roll back (no applied migrations)")
			return
		}
		fmt.Fprintf(cmd.OutOrStdout(), "rollback OK: %s\n", group.String())
	},
}

func init() {
	rootCmd.AddCommand(rollbackCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rollbackCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rollbackCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
