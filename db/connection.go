package db

import (
	"context"
	"fmt"
	"github.com/Gohelraj/youtube-search-api/config"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"runtime"
)

func connectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Conf.Database.Host, config.Conf.Database.Port, config.Conf.Database.User, config.Conf.Database.Password, config.Conf.Database.Name, config.Conf.Database.SSLMode)
}

// Connect creates a connection to the database
func Connect() (*pgxpool.Pool, error) {
	conf, err := pgxpool.ParseConfig(connectionString())
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