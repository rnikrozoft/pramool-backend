package cmd

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

// ทุก migration (core / wallet / auction) รันลง database เดียวตาม DATABASE_NAME
func openMigrateDB(_ string) (*bun.DB, string, error) {
	dsn, err := migrateDSN()
	if err != nil {
		return nil, "", err
	}
	label := redactDSN(dsn)
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	return bun.NewDB(sqldb, pgdialect.New()), label, nil
}

func migrateDSN() (string, error) {
	if dsn := strings.TrimSpace(os.Getenv("DATABASE_DSN")); dsn != "" {
		return dsn, nil
	}
	if dsn := dsnFromParts(os.Getenv("DATABASE_NAME")); dsn != "" {
		return dsn, nil
	}
	return "", fmt.Errorf("set DATABASE_DSN or DATABASE_HOST + DATABASE_USERNAME + DATABASE_NAME")
}

func dsnFromParts(dbName string) string {
	if strings.TrimSpace(dbName) == "" {
		return ""
	}
	host := strings.TrimSpace(os.Getenv("DATABASE_HOST"))
	user := strings.TrimSpace(os.Getenv("DATABASE_USERNAME"))
	pass := os.Getenv("DATABASE_PASSWORD")
	port := strings.TrimSpace(os.Getenv("DATABASE_PORT"))
	if host == "" || user == "" {
		return ""
	}
	if port == "" {
		port = "5432"
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, dbName)
}

func redactDSN(dsn string) string {
	schemeEnd := strings.Index(dsn, "://")
	if schemeEnd < 0 {
		return dsn
	}
	rest := dsn[schemeEnd+3:]
	at := strings.Index(rest, "@")
	if at < 0 {
		return dsn
	}
	userPart := rest[:at]
	if colon := strings.Index(userPart, ":"); colon >= 0 {
		userPart = userPart[:colon]
	}
	return dsn[:schemeEnd+3] + userPart + rest[at:]
}
