package db

import (
	"context"
	"fmt"
	"os"
	"testing"
)

func TestDBConnection(t *testing.T) {
	pool := MountDB()
	defer pool.Close()

	// Define test data
	_, err := pool.Exec(context.Background(), "SELECT 1")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ping failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Ping successful")
}
