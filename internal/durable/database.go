package durable

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Struct for database connection info
type ConnectionInfo struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

// Struct for database connection
type Database struct {
	db *pgxpool.Pool
}

// OpenDatabaseClient open a new database connection
func OpenDatabaseClient(ctx context.Context, c *ConnectionInfo) *pgxpool.Pool {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", c.User, c.Password, c.Host, c.Port, c.Name)

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Panicln(err)
	}
	config.MinConns = 10
	config.MaxConns = 10

	dppool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Panicln(err)
	}

	return dppool
}

// WrapDatabase put database connection into Database struct
func WrapDatabase(db *pgxpool.Pool) *Database {
	return &Database{db}
}

// Close close current database connection
func (d *Database) Close() {
	d.db.Close()
}

// Exec execute sql query
func (d *Database) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	return d.db.Exec(ctx, query, args...)
}

// Query execute sql query with the given arguments
func (d *Database) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	return d.db.Query(ctx, query, args...)
}

// QueryRow execute sql query with the given arguments
func (d *Database) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return d.db.QueryRow(ctx, query, args...)
}
