package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgres() *pgxpool.Pool {
    url := os.Getenv("DATABASE_URL")
    if url == "" {
        url = "postgres://postgres:postgres@db:5432/rules"
    }

    pool, err := pgxpool.New(context.Background(), url)
    if err != nil {
        log.Fatal(err)
    }

    return pool
}