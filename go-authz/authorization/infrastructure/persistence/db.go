package persistence

import (
	"authorization/config"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	DBConnection *pgxpool.Pool
)

func CreateDBConnection() error {
	// newLogger := logger.New(
	// 	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
	// 	logger.Config{
	// 		SlowThreshold:             time.Second, // Slow SQL threshold
	// 		LogLevel:                  logger.Info, // Log level
	// 		IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
	// 		Colorful:                  false,       // Disable color
	// 	},
	// )

	dsn := fmt.Sprintf("%s://%s:%s@%s:%s/%s", config.StorageConfig.DBDriver, config.StorageConfig.DBUser, config.StorageConfig.DBPassword, config.StorageConfig.DBHost, config.StorageConfig.DBPort, config.StorageConfig.DBName)
	ctx := context.Background()
	var err error

	DBConnection, err = pgxpool.New(ctx, dsn)

	if err != nil {
		return err
	}

	// connection, err := gorm.Open(postgres.New(postgres.Config{
	// 	DriverName: "pgx",
	// 	DSN:        dsn,
	// }), &gorm.Config{
	// 	PrepareStmt: true,
	// 	Logger:      newLogger,
	// })
	// if err != nil {
	// 	return err
	// }

	// pgDB, err := connection.DB()
	// if err != nil {
	// 	return err
	// }

	// // SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	// pgDB.SetMaxIdleConns(config.StorageConfig.MaxIdleConns)

	// // SetMaxOpenConns sets the maximum number of open connections to the database.
	// pgDB.SetMaxOpenConns(config.StorageConfig.MaxOpenConns)

	// // SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	// duration := time.Duration(config.StorageConfig.ConnMaxLifetime * int64(time.Minute))
	// pgDB.SetConnMaxLifetime(duration)

	return nil
}
