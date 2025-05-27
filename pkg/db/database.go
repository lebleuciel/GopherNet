package db

import (
	"context"
	"fmt"
	"time"

	dbsql "database/sql"
	"gophernet/pkg/config"
	"gophernet/pkg/db/ent"
	"gophernet/pkg/shutdown"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

type Database interface {
	Close() error
	EntClient() *ent.Client
	DB() *dbsql.DB
	IsInitialized(ctx context.Context) (bool, error)
}

type database struct {
	pool     *pgxpool.Pool
	database *dbsql.DB
	client   *ent.Client
}

func NewDatabase(ctx context.Context, dbConfig *config.Database) Database {
	if dbConfig == nil {
		panic("database config cannot be nil")
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Database,
	)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		panic(err)
	}

	// Configure pool
	poolConfig.MaxConns = 10
	poolConfig.MinConns = 2

	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(timeoutCtx, poolConfig)
	if err != nil {
		panic(err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		panic(err)
	}

	db := stdlib.OpenDB(*pool.Config().ConnConfig)

	// Create ent driver
	driver := sql.OpenDB(dialect.Postgres, db)

	// Create ent client
	client := ent.NewClient(ent.Driver(driver))

	// Run automatic migrations
	if err := client.Schema.Create(ctx); err != nil {
		panic(fmt.Errorf("failed creating schema resources: %w", err))
	}

	// Register shutdown handler
	shutdown.GetManager().Register("database", func(ctx context.Context) error {
		return client.Close()
	})

	return &database{
		pool:     pool,
		database: db,
		client:   client,
	}
}

func (db *database) Close() error {
	var errs []error

	if err := db.client.Close(); err != nil {
		errs = append(errs, fmt.Errorf("closing ent client: %w", err))
	}

	db.pool.Close()
	if err := db.database.Close(); err != nil {
		errs = append(errs, fmt.Errorf("closing database: %w", err))
	}

	if len(errs) > 0 {
		return fmt.Errorf("closing database: %v", errs)
	}
	return nil
}

func (db *database) EntClient() *ent.Client {
	return db.client
}

func (db *database) DB() *dbsql.DB {
	return db.database
}

// IsInitialized checks if the database is properly initialized
func (d *database) IsInitialized(ctx context.Context) (bool, error) {
	// Check if the burrows table exists
	var exists bool
	err := d.database.QueryRow(`
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'burrows'
		);
	`).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if burrows table exists: %w", err)
	}

	return exists, nil
}
