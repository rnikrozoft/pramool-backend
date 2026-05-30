/*
Copyright © 2025 rnikrozoft rnikrozoft.dev@gmail.com
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rnikrozoft/pramool-core/migrations"
	"github.com/spf13/cobra"
)

// newMigrationCmd creates a new .up.sql / .down.sql pair under migrations/<db>/.
var newMigrationCmd = &cobra.Command{
	Use:   "newMigration <db> <name>",
	Short: "Create a new SQL migration pair under migrations/<db>/",
	Long: `db is one of: core, wallet, auction

Example:
  go run . newMigration wallet my_table`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		db := strings.TrimSpace(strings.ToLower(args[0]))
		name := strings.TrimSpace(args[1])
		switch db {
		case migrations.DBCore, migrations.DBWallet, migrations.DBAuction:
		default:
			return errors.New("db must be core, wallet, or auction")
		}
		if name == "" {
			return errors.New("migration name is required")
		}
		stamp := time.Now().UTC().Format("20060102150405")
		base := fmt.Sprintf("%s_%s", stamp, name)
		dir := filepath.Join("migrations", db)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
		body := []byte("SET statement_timeout = 0;\n\n--bun:split\n\n")
		upPath := filepath.Join(dir, base+".up.sql")
		downPath := filepath.Join(dir, base+".down.sql")
		if err := os.WriteFile(upPath, body, 0o644); err != nil {
			return err
		}
		if err := os.WriteFile(downPath, body, 0o644); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "created %s\n", upPath)
		fmt.Fprintf(cmd.OutOrStdout(), "created %s\n", downPath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(newMigrationCmd)
}
