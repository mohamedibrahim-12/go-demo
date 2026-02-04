package database

import (
	"database/sql"
	"fmt"
	"os"

	"go-demo/pkg/logger"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	// DSN constructed from env vars

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("failed to open DB")
	}

	if err = DB.Ping(); err != nil {
		logger.Log.Fatal().Err(err).Msg("failed to ping DB")
	}

	logger.Log.Info().Str("db", os.Getenv("DB_NAME")).Msg("connected to PostgreSQL")

	// Run lightweight migrations to ensure `uuid` columns exist and pgcrypto is available.
	migrations := []string{
		`CREATE EXTENSION IF NOT EXISTS "pgcrypto";`,
		`CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, uuid UUID NOT NULL DEFAULT gen_random_uuid(), name TEXT NOT NULL, role TEXT NOT NULL);`,
		`CREATE TABLE IF NOT EXISTS products (id SERIAL PRIMARY KEY, uuid UUID NOT NULL DEFAULT gen_random_uuid(), name TEXT NOT NULL, price NUMERIC NOT NULL);`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS uuid UUID;`,
		`UPDATE users SET uuid = gen_random_uuid() WHERE uuid IS NULL;`,
		`ALTER TABLE users ALTER COLUMN uuid SET NOT NULL;`,
		`ALTER TABLE users ALTER COLUMN uuid SET DEFAULT gen_random_uuid();`,
		`ALTER TABLE products ADD COLUMN IF NOT EXISTS uuid UUID;`,
		`UPDATE products SET uuid = gen_random_uuid() WHERE uuid IS NULL;`,
		`ALTER TABLE products ALTER COLUMN uuid SET NOT NULL;`,
		`ALTER TABLE products ALTER COLUMN uuid SET DEFAULT gen_random_uuid();`,
	}

	for _, m := range migrations {
		if _, err := DB.Exec(m); err != nil {
			logger.Log.Error().Err(err).Str("migration", m).Msg("migration statement failed")
		}
	}
}
