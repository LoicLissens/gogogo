package db

import (
	"context"
	"fmt"
	"os"
	"testing"

	"jiva-guildes/settings"
)

func TestDBConnection(t *testing.T) {
	pool := MountDB(settings.AppSettings.DATABASE_URI)
	defer pool.Close()

	_, err := pool.Exec(context.Background(), "SELECT 1")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ping failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Ping successful")
}
