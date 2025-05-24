/*
Copyright © 2025 rnikrozoft rnikrozoft.dev@gmail.com
*/
package cmd

import (
	"context"
	"errors"

	"github.com/rnikrozoft/pramool.in.th-backend/migrations"
	"github.com/spf13/cobra"
)

// newMigrationCmd represents the newMigration command
var newMigrationCmd = &cobra.Command{
	Use:   "newMigration",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PreRunE: func(_ *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("newMigration command need one argument (file name)")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		migrator, err := migrations.GetMigrator(context.Background(), nil)
		if err != nil {
			panic(err)
		}
		migrator.CreateSQLMigrations(context.Background(), args[0])
	},
}

func init() {
	rootCmd.AddCommand(newMigrationCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newMigrationCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newMigrationCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
