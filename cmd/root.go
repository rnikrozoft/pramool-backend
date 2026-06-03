/*
Copyright © 2025 rnikrozoft rnikrozoft.dev@gmail.com
*/
package cmd

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/rnikrozoft/pramool-core/config"
	"github.com/spf13/cobra"
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
	// Keep behavior consistent with other services: load .env if present, but do not require it.
	_ = godotenv.Load()

	appConfigs = config.AppConfigs{
		Database: config.DatabaseConfig{
			Host:         os.Getenv("DATABASE_HOST"),
			Port:         os.Getenv("DATABASE_PORT"),
			Username:     os.Getenv("DATABASE_USERNAME"),
			Password:     os.Getenv("DATABASE_PASSWORD"),
			DatabaseName: os.Getenv("DATABASE_NAME"),
		},
		Jwt: config.JwtConfig{
			Issuer:            os.Getenv("JWT_ISSUER"),
			Secret:            os.Getenv("JWT_SECRET"),
			ExpireTime:        envInt("JWT_EXPIRE_TIME", 0),
			RefreshExpireTime: envInt("JWT_REFRESH_EXPIRE_TIME", 0),
		},
		ThaiBulkSMS: config.ThaiBulkSMS{
			AddressRequest: os.Getenv("ADDRESS_REQUEST"),
			AddressVerify:  os.Getenv("ADDRESS_VERIFY"),
			APIKey:         os.Getenv("API_KEY"),
			APISecret:      os.Getenv("API_SECRET"),
		},
		CorsAllowOrigins:   os.Getenv("CORS_ALLOW_ORIGINS"),
		TrackingMoreAPIKey: os.Getenv("TRACKINGMORE_API_KEY"),
		NationalIDEncKey:   os.Getenv("NATIONAL_ID_ENCRYPTION_KEY"),
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

func envInt(name string, fallback int) int {
	v := os.Getenv(name)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}

// databaseTargetLine returns a password-redacted DSN label for CLI output (migrate / rollback).
func databaseTargetLine() string {
	d := appConfigs.Database
	return fmt.Sprintf("postgres://%s@%s:%s/%s", d.Username, d.Host, d.Port, d.DatabaseName)
}
