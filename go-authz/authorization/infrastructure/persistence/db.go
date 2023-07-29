package persistence

import (
	"authorization/config"
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"github.com/golang-migrate/migrate/v4"
	pgxadapter "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var (
	Pool *pgxpool.Pool
)

func CreateDBConnection() error {
	dsn := fmt.Sprintf("%s://%s:%s@%s:%s/%s", config.StorageConfig.DBDriver, config.StorageConfig.DBUser, config.StorageConfig.DBPassword, config.StorageConfig.DBHost, config.StorageConfig.DBPort, config.StorageConfig.DBName)
	ctx := context.Background()
	var err error

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return err
	}

	Pool, err = pgxpool.NewWithConfig(ctx, config)

	if err != nil {
		return err
	}

	return nil
}

func Migration(newPool *pgxpool.Pool) {
	if newPool != nil {
		Pool = newPool
	}

	conn, err := sql.Open("pgx", Pool.Config().ConnString())
	if err != nil {
		log.Fatal().Caller().Err(err).Msg("Failed to connect to database")
	}

	driver, err := pgxadapter.WithInstance(conn, &pgxadapter.Config{})
	if err != nil {
		log.Fatal().Caller().Err(err).Msg("Failed to create migration driver")
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal().Caller().Err(err).Msg("Failed to get working directory")
	}

	path := fmt.Sprintf("%s%s%s", "file://", wd, "/infrastructure/persistence/migrations")
	m, err := migrate.NewWithDatabaseInstance(path, "pgx", driver)
	if err != nil {
		log.Fatal().Caller().Err(err).Msg("Failed to find migration files")
	}

	err = m.Up()
	if err != nil {
		if err.Error() != "no change" {
			log.Fatal().Caller().Err(err).Msg("Failed to migrate database")
		}
		log.Info().Caller().Msg("No database changes")
	} else {
		log.Info().Caller().Msg("Migrated database")

	}

	conn.Close() // close sql connection since we use pgxpool pgxv5
}
