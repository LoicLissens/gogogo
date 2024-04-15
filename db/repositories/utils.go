package repositories

import (
	"context"
	"fmt"
	"jiva-guildes/db/tables"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	GetByUUID(connectionPool *pgxpool.Pool, uuid uuid.UUID, tableName string, schema string) (interface{}, error)
	Save(connectionPool *pgxpool.Pool, tableName string, schema string, entity interface{}) (interface{}, error)
}

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
