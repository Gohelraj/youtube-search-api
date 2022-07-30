package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Connect creates a connection to the database
func Connect(dbURL string) *pgxpool.Pool {
	pool, err := pgxpool.Connect(context.Background(), dbURL)
	if err != nil {
		log.Println(err)
	}
	if err := pool.Ping(context.TODO()); err != nil {
		log.Fatal(err)
	}
	return pool
}