package di

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/KauanCarvalho/fiap-sa-product-service/internal/config"

	_ "github.com/go-sql-driver/mysql" //nolint:riveve // its mandatory.
	"github.com/go-testfixtures/testfixtures/v3"
	_ "github.com/golang-migrate/migrate/v4/source/file" //nolint:riveve // its mandatory.
	gormMySQL "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	connMaxLifetime          = 5 * time.Minute
	defaultConnectionTimeout = 5 * time.Second
)

func NewDatabaseConnectionPool(cfg *config.Config) (*gorm.DB, error) {
	var logMode logger.LogLevel
	if cfg.IsProduction() {
		logMode = logger.Silent
	} else {
		logMode = logger.Info
	}

	db, errOpen := gorm.Open(gormMySQL.Open(formatDSN(cfg)), &gorm.Config{
		Logger: logger.Default.LogMode(logMode),
	})
	if errOpen != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", errOpen)
	}

	sqlDB, errDB := db.DB()
	if errDB != nil {
		return nil, errDB
	}

	sqlDB.SetConnMaxLifetime(connMaxLifetime)
	sqlDB.SetMaxOpenConns(cfg.DatabaseMaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.DatabaseMaxIdleConns)

	ctx, cancel := context.WithTimeout(context.Background(), defaultConnectionTimeout)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func formatDSN(cfg *config.Config) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true",
		cfg.DatabaseUser,
		cfg.DatabasePassword,
		cfg.DatabaseHost,
		cfg.DatabasePort,
		cfg.DatabaseName,
	)
}

func SetupFixtures(db *sql.DB, path string) (*testfixtures.Loader, error) {
	return testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("mysql"),
		testfixtures.Directory(path),
	)
}
