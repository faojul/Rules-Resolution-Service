package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(url string) *pgxpool.Pool {
	if url == "" {
		url = "postgres://spine:spine@postgres:5432/rules_resolution?sslmode=disable"
	}

	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		log.Fatal(err)
	}

	return pool
}