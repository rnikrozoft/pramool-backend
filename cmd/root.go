/*
Copyright © 2025 rnikrozoft rnikrozoft.dev@gmail.com
*/
package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/rnikrozoft/pramool-core/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var rootCmd = &cobra.Command{
	Use:   "pramool-core",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

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
			Issuer:            viper.GetString("JWT_ISSUER"),
			Secret:            viper.GetString("JWT_SECRET"),
			ExpireTime:        viper.GetInt("JWT_EXPIRE_TIME"),
			RefreshExpireTime: viper.GetInt("JWT_REFRESH_EXPIRE_TIME"),
		},
		ThaiBulkSMS: config.ThaiBulkSMS{
			AddressRequest: viper.GetString("ADDRESS_REQUEST"),
			AddressVerify:  viper.GetString("ADDRESS_VERIFY"),
			APIKey:         viper.GetString("API_KEY"),
			APISecret:      viper.GetString("API_SECRET"),
		},
		CorsAllowOrigins: viper.GetString("CORS_ALLOW_ORIGINS"),
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

// databaseTargetLine returns a password-redacted DSN label for CLI output (migrate / rollback).
func databaseTargetLine() string {
	d := appConfigs.Database
	return fmt.Sprintf("postgres://%s@%s:%s/%s", d.Username, d.Host, d.Port, d.DatabaseName)
}
