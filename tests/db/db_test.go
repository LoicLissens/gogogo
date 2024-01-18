package test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"jiva-guildes/db"
	"jiva-guildes/settings"
)

var schema string = settings.AppSettings.DATABASE_SCHEMA

func TestDBConnection(t *testing.T) {
	pool := db.MountDB()
	defer pool.Close()

	// Define test data	var schema string = settings.AppSettings.DATABASE_SCHEMA
	_, err := pool.Exec(context.Background(), "SELECT 1")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ping failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Ping successful")
}
