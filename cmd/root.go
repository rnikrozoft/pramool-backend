/*
Copyright © 2025 rnikrozoft rnikrozoft.dev@gmail.com
*/
package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/rnikrozoft/pramool.in.th-backend/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pramool.in.th-backend",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var conn *bun.DB
var appConfigs config.AppConfigs

func init() {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	appConfigs = config.AppConfigs{
		Database: config.DatabaseConfig{
			Host:         viper.GetString("DATABASE_HOST"),
			Port:         viper.GetString("DATABASE_PORT"),
			Username:     viper.GetString("DATABASE_USERNAME"),
			Password:     viper.GetString("DATABASE_PASSWORD"),
			DatabaseName: viper.GetString("DATABASE_NAME"),
		},
		Jwt: config.JwtConfig{
			Secret:     viper.GetString("JWT_SECRET"),
			ExpireTime: viper.GetInt("JWT_EXPIRE_TIME"),
		},
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		appConfigs.Database.Username,
		appConfigs.Database.Password,
		appConfigs.Database.Host,
		appConfigs.Database.Port,
		appConfigs.Database.DatabaseName)

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	conn = bun.NewDB(sqldb, pgdialect.New())
}
