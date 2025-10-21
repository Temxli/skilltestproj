package database

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ...existing code...
var DB *gorm.DB

func DBConnect() error {
	// load local .env if present
	_ = godotenv.Load()

	dsn := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if dsn == "" {
		host := trim(firstNonEmpty(os.Getenv("DB_HOST"), "localhost"))
		port := trim(firstNonEmpty(os.Getenv("DB_PORT"), "5432"))
		user := trim(os.Getenv("DB_USER"))
		password := trim(os.Getenv("DB_PASSWORD"))
		name := trim(os.Getenv("DB_NAME"))
		sslmode := trim(firstNonEmpty(os.Getenv("DB_SSLMODE"), "disable"))

		if user == "" || password == "" || name == "" {
			return fmt.Errorf("missing required env vars DB_USER, DB_PASSWORD or DB_NAME")
		}

		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			host, user, password, name, port, sslmode)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	DB = db
	return nil
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

func trim(s string) string {
	return strings.TrimSpace(s)
}

// ...existing code...
