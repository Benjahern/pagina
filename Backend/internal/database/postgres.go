package database

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"tienda/backend/internal/config"
)

// Connection pool defaults. Tuned for a small-to-medium API server.
// Move to config/env if traffic patterns demand different values.
const (
	defaultMaxOpenConns    = 25
	defaultMaxIdleConns    = 5
	defaultConnMaxLifetime = 5 * time.Minute
	defaultConnMaxIdleTime = 10 * time.Minute
	defaultPingTimeout     = 5 * time.Second
)

// Connect opens a GORM connection to PostgreSQL, configures the pool, and
// verifies reachability with a ping. Returns a ready-to-use *gorm.DB.
//
// SECURITY — defense against SQL injection:
// GORM parameterizes queries by default. Where("email = ?", input) sends
// `input` as a bound parameter ($1), never as interpolated SQL — so a value
// like "x'; DROP TABLE users; --" is treated as data, not code. This is the
// primary defense at the data layer.
//
// DO use:
//   db.Where("email = ?", userInput).First(&u)
//   db.Raw("SELECT * FROM users WHERE id = ?", id).Scan(&u)
//
// DO NOT build SQL with fmt.Sprintf/user input — that re-opens the injection
// vector GORM is closing for us:
//   db.Where(fmt.Sprintf("email = '%s'", userInput)) // BAD
//
// The DSN string here only carries connection params (host/port/user/pass/db)
// — it is never used as a query body, so building it via Sprintf is safe.
//
// sslmode=disable is acceptable for the local docker container; production
// should set sslmode=require or verify-full via config.
func Connect(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("database: open postgres at %s:%d: %w", cfg.Host, cfg.Port, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("database: get underlying *sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(defaultMaxOpenConns)
	sqlDB.SetMaxIdleConns(defaultMaxIdleConns)
	sqlDB.SetConnMaxLifetime(defaultConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(defaultConnMaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), defaultPingTimeout)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("database: ping postgres at %s:%d: %w", cfg.Host, cfg.Port, err)
	}

	return db, nil
}