package database

import (
	"database/sql"
	"fmt"
	"os"

	"go-demo/models"
	"go-demo/pkg/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
)

var (
	DB     *sql.DB
	GormDB *gorm.DB
)

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

	// Initialize GORM on top of the DSN
	var err error
	GormDB, err = gorm.Open(postgres.New(postgres.Config{DSN: dsn, PreferSimpleProtocol: true}), &gorm.Config{})
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("failed to open GORM DB")
	}

	// get underlying sql.DB from GORM and keep for compatibility
	sqlDB, err := GormDB.DB()
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("failed to get underlying sql.DB from Gorm")
	}
	DB = sqlDB

	if err = DB.Ping(); err != nil {
		logger.Log.Fatal().Err(err).Msg("failed to ping DB")
	}

	// Migration gating using a lightweight schema_migrations table.
	// This avoids per-column checks when models grow large.
	migr := GormDB.Migrator()

	// Ensure the migrations table exists (simple single-row key table)
	if !migr.HasTable("schema_migrations") {
		if err := GormDB.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (version text PRIMARY KEY, applied_at timestamptz DEFAULT now())`).Error; err != nil {
			logger.Log.Warn().Err(err).Msg("failed to create schema_migrations table; continuing")
		}
	}

	// Check whether our auto-migrate version has already been applied
	var appliedCount int64
	const migrationVersion = "auto_migrate_v1"
	if err := GormDB.Raw(`SELECT COUNT(1) FROM schema_migrations WHERE version = ?`, migrationVersion).Scan(&appliedCount).Error; err != nil {
		// If the query fails for unexpected reasons, fall back to attempting migration
		logger.Log.Warn().Err(err).Msg("failed to query schema_migrations; will attempt AutoMigrate")
		appliedCount = 0
	}

	if appliedCount == 0 {
		// Ensure pgcrypto extension exists (needed for gen_random_uuid())
		if res := GormDB.Exec(`CREATE EXTENSION IF NOT EXISTS "pgcrypto";`); res.Error != nil {
			logger.Log.Warn().Err(res.Error).Msg("failed to ensure pgcrypto extension; continuing and hoping extension exists")
		}

		if err := GormDB.AutoMigrate(&models.User{}, &models.Product{}, &models.AuditLog{}, &models.NotificationJob{}); err != nil {
			logger.Log.Fatal().Err(err).Msg("auto-migrate failed")
		}

		if err := GormDB.Exec(`INSERT INTO schema_migrations (version) VALUES (?) ON CONFLICT DO NOTHING`, migrationVersion).Error; err != nil {
			logger.Log.Warn().Err(err).Msg("failed to record applied migration version; migration still applied")
		}

		logger.Log.Info().Str("migration", migrationVersion).Msg("auto-migrate completed and recorded")
	} else {
		logger.Log.Info().Str("migration", migrationVersion).Msg("migration already applied; skipping AutoMigrate")
	}

	// v2: ensure created_at columns exist for cleanup worker (idempotent)
	for _, q := range []string{
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()`,
		`ALTER TABLE products ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()`,
	} {
		if err := GormDB.Exec(q).Error; err != nil {
			logger.Log.Fatal().Err(err).Msg("auto-migrate v2 (created_at) failed")
		}
	}
	const migrationV2 = "auto_migrate_v2"
	_ = GormDB.Exec(`INSERT INTO schema_migrations (version) VALUES (?) ON CONFLICT DO NOTHING`, migrationV2).Error

	// v3: ensure worker tables exist
	if err := GormDB.AutoMigrate(&models.AuditLog{}, &models.NotificationJob{}); err != nil {
		logger.Log.Fatal().Err(err).Msg("auto-migrate v3 (workers) failed")
	}
	const migrationV3 = "auto_migrate_v3"
	_ = GormDB.Exec(`INSERT INTO schema_migrations (version) VALUES (?) ON CONFLICT DO NOTHING`, migrationV3).Error

	logger.Log.Info().Str("db", os.Getenv("DB_NAME")).Msg("connected to PostgreSQL")
}
