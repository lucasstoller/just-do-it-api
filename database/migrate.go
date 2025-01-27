package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func createMigrator(sqlDB *sql.DB) (*migrate.Migrate, error) {
	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("could not create the postgres driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create migrate instance: %v", err)
	}

	return m, nil
}

// RunMigrations executes all pending migrations
func RunMigrations(dsn string) error {
	gormDB := initDB()
	sqlDB, err := gormDB.DB()
	if err != nil {
		return fmt.Errorf("could not get underlying *sql.DB: %v", err)
	}

	// Create schema_migrations table first
	if _, err := sqlDB.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version bigint NOT NULL,
			dirty boolean NOT NULL,
			CONSTRAINT schema_migrations_pkey PRIMARY KEY (version)
		);
	`); err != nil {
		return fmt.Errorf("could not create schema_migrations table: %v", err)
	}

	m, err := createMigrator(sqlDB)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not run migrations: %v", err)
	}

	log.Println("Migrations completed successfully")
	return nil
}

// ResetDatabase drops all tables and reruns all migrations
func ResetDatabase(dsn string) error {
	gormDB := initDB()
	sqlDB, err := gormDB.DB()
	if err != nil {
		return fmt.Errorf("could not get underlying *sql.DB: %v", err)
	}

	// Drop all tables
	if _, err := sqlDB.Exec(`
		DROP TABLE IF EXISTS tasks CASCADE;
		DROP TABLE IF EXISTS users CASCADE;
		DROP TABLE IF EXISTS schema_migrations CASCADE;
	`); err != nil {
		return fmt.Errorf("could not drop tables: %v", err)
	}

	// Create schema_migrations table
	if _, err := sqlDB.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version bigint NOT NULL,
			dirty boolean NOT NULL,
			CONSTRAINT schema_migrations_pkey PRIMARY KEY (version)
		);
	`); err != nil {
		return fmt.Errorf("could not create schema_migrations table: %v", err)
	}

	m, err := createMigrator(sqlDB)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not run migrations: %v", err)
	}

	log.Println("Database reset completed successfully")
	return nil
}
