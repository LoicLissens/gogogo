package db

import (
	"context"
	"fmt"
	"jiva-guildes/settings"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func MountDB() *pgxpool.Pool {
	dbpool, err := pgxpool.New(context.Background(), settings.AppSettings.DATABASE_URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return dbpool
}

func Teardown(connectionPool *pgxpool.Pool) {
	connectionPool.Close()
}
