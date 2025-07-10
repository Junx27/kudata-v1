package database

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"payment/internal/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(cfg config.Config) error {
	pass := url.QueryEscape(cfg.DatabasePassword)
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		cfg.DatabaseUsername,
		pass,
		fmt.Sprintf("%s:%s", cfg.DatabaseHost, cfg.DatabasePort),
		cfg.DatabaseName,
	)

	// Open DB for postgres driver
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Printf("Failed to open DB: %v\n", err)
		return err
	}

	// Initialize driver
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Printf("Failed to create postgres driver: %v\n", err)
		return err
	}

	// Create migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", cfg.MigrationPath),
		"postgres", driver)
	if err != nil {
		log.Printf("Failed to create migrate instance: %v\n", err)
		return err
	}

	// Check version and dirty status
	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.Printf("Failed to get current migration version: %v\n", err)
		return err
	}

	if dirty {
		log.Printf("Database is dirty at version %d. Forcing clean...\n", version)
		err = m.Force(int(version))
		if err != nil {
			log.Printf("Failed to force clean: %v\n", err)
			return err
		}
	}

	// Run migrations
	log.Println("Running database migrations...")
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Printf("Migration failed: %v\n", err)
		return fmt.Errorf("migration failed: %v", err)
	}

	log.Println("Database migration completed successfully.")
	return nil
}
