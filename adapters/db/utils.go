package db

import (
	"errors"
	"fmt"
	customerrors "jiva-guildes/domain/custom_errors"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func HandleSQLErrors(err error, tableName string, uuid uuid.UUID) error {
	if errors.Is(err, pgx.ErrNoRows) {
		errorMessage := fmt.Sprintf("No entity with UUID %s found in table %s", uuid, tableName)
		return customerrors.NewErrorNotFound(errorMessage)
	}
	if e, ok := err.(*pgconn.PgError); ok && e.Code == pgerrcode.UniqueViolation {
		errorMessage := fmt.Sprintf("Entity with UUID %s already exists in table %s", uuid, tableName)
		return customerrors.NewErrorAlreadyExists(errorMessage)
	}
	if e, ok := err.(*pgconn.PgError); ok {
		fmt.Println(e.Code)
	}
	log.Fatal(err)
	return err
}
