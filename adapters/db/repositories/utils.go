package repositories

import (
	"context"
	"fmt"
	"jiva-guildes/adapters/db/tables"
	customerrors "jiva-guildes/domain/custom_errors"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetEntityByUuid(connectionPool *pgxpool.Pool, uuid uuid.UUID, tableName string) pgx.Row {
	statement := fmt.Sprintf("SELECT * FROM %s WHERE uuid = $1", tableName)
	row := connectionPool.QueryRow(context.Background(), statement, uuid)

	return row
}

func SaveEntity(table tables.Table, conn *pgxpool.Pool) pgx.Row {
	tableFields, values := tables.DeepFields(table)
	fields := ""
	fieldsPosition := ""
	for i, field := range tableFields {
		if i != 0 {
			fields += ", "
			fieldsPosition += ", "
		}
		fields += fmt.Sprintf("%s", tables.GetDBColumnName(field.Name, table))
		fieldsPosition += fmt.Sprintf("$%d", i+1)
	}
	statement := fmt.Sprintf(`INSERT INTO %s(%s) VALUES(%s) RETURNING *;`, table.GetTableName(), fields, fieldsPosition)
	interfaceValues := make([]interface{}, len(values))
	for i, v := range values {
		interfaceValues[i] = v.Interface()
	}
	row := conn.QueryRow(context.Background(), statement, interfaceValues...)
	return row
}
func DeleteEntity(tableName string, uuid uuid.UUID, conn *pgxpool.Pool) (int64, error) {
	statement := fmt.Sprintf("DELETE FROM %s WHERE uuid = $1", tableName)
	result, err := conn.Exec(context.Background(), statement, uuid)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), err
}

func HandleSQLDelete(rowAffected int64, err error, tableName string, uuid uuid.UUID) error {
	if err != nil {
		log.Fatal(err)
		return err
	}
	if rowAffected == 0 {
		errorMessage := fmt.Sprintf("No entity with UUID %s found in table %s", uuid, tableName)
		return customerrors.NewErrorNotFound(errorMessage)
	}
	return nil
}
