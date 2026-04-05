package cmd

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // register postgres driver
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/spf13/cobra"

	"github.com/example/grinex-rates-service/internal/config"
	migrations "github.com/example/grinex-rates-service/internal/migration/postgres"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply all pending migrations",
	RunE: func(cmd *cobra.Command, _ []string) error {
		m, err := newMigrate(cmd)
		if err != nil {
			return err
		}
		defer func() { _, _ = m.Close() }()

		if err := m.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("no pending migrations")
				return nil
			}

			return fmt.Errorf("migrate up: %w", err)
		}

		fmt.Println("migrations applied")
		return nil
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Rollback the last migration",
	RunE: func(cmd *cobra.Command, _ []string) error {
		m, err := newMigrate(cmd)
		if err != nil {
			return err
		}
		defer func() { _, _ = m.Close() }()

		if err := m.Steps(-1); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("no migrations to roll back")
				return nil
			}

			return fmt.Errorf("migrate down: %w", err)
		}

		fmt.Println("migration rolled back")
		return nil
	},
}

func init() {
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
}

func newMigrate(cmd *cobra.Command) (*migrate.Migrate, error) {
	cfg, err := config.Load(cmd)
	if err != nil {
		return nil, err
	}

	if cfg.Postgres.DSN == "" {
		return nil, fmt.Errorf("database-dsn is required for migrations")
	}

	source, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return nil, fmt.Errorf("migration source: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", source, cfg.Postgres.DSN)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	return m, nil
}
