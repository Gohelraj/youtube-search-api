package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"runtime"
)

// Connect creates a connection to the database
func Connect(dbURL string) (*pgxpool.Pool, error) {
	conf, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, err
	}

	// pgxpool default max number of connections is the number of CPUs on your machine returned by runtime.NumCPU().
	// This number is very conservative, and you might be able to improve performance for highly concurrent applications
	// by increasing it.
	conf.MaxConns = int32(runtime.NumCPU() * 2)

	pool, err := pgxpool.ConnectConfig(context.Background(), conf)
	if err != nil {
		return nil, fmt.Errorf("pgx connection error: %w", err)
	}
	if err := pool.Ping(context.TODO()); err != nil {
		log.Fatal(err)
	}
	return pool, nil
}