package db

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// To instanciate with  settings.AppSettings.DATABASE_URI
func MountDB(databaseURI *url.URL) *pgxpool.Pool {
	dbpool, err := pgxpool.New(context.Background(), databaseURI.Path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return dbpool
}

func Teardown(connectionPool *pgxpool.Pool) {
	connectionPool.Close()
}
