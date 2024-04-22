package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func MountDB(databaseURI string) *pgxpool.Pool {
	dbpool, err := pgxpool.New(context.Background(), databaseURI)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return dbpool
}

func Teardown(connectionPool *pgxpool.Pool) {
	connectionPool.Close()
}
