package repositories

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetEntityByUuid(connectionPool *pgxpool.Pool, uuid uuid.UUID, tableName string, schema string) (interface{}, error) {
	statement := `SELECT * FROM $1.$2 WHERE uuid = $3`
	row := connectionPool.QueryRow(context.Background(), statement, schema, tableName, uuid)

	var results interface{}
	err := row.Scan(&results)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("no entity with UUID %s found in table %s.%s", uuid, schema, tableName)
		} else {
			return nil, err
		}
	}
	return results, nil
}

func SaveEntity(connectionPool *pgxpool.Pool, tableName string, schema string, entity interface{}) (interface{}, error) {
	var valuesParams strings.Builder

	entityType := reflect.TypeOf(entity)
	entityValue := reflect.ValueOf(entity)
	numFields := entityType.NumField()
	values := make([]interface{}, numFields)

	for i := 0; i < numFields; i++ {
		field := entityType.Field(i)
		valuesParams.WriteString(fmt.Sprintf("$%v,", i+1)) // Values should be represent as $1,$2 to prevent sql injection.
		values[i] = entityValue.FieldByName(field.Name).Interface()
	}
	prams := strings.TrimRight(valuesParams.String(), ",")
	statement := fmt.Sprintf(`INSERT INTO %s.%s VALUES (%s) RETURNING *`, schema, tableName, prams)

	rows, err := connectionPool.Query(context.Background(), statement, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// fields := rows.FieldDescriptions()
	// savedValues := make([]interface{}, len(fields))
	// valuePtrs := make([]interface{}, len(fields))
	// for i := range savedValues {
	// 	valuePtrs[i] = &savedValues[i]
	// }

	// for rows.Next() {
	// 	err = rows.Scan(valuePtrs...)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	// Now `values` contains the values from the first row, and you can process them as needed.
	// If the query could return multiple rows, you would need to handle that as well.
	lol := reflect.New(entityType)

	// for rows.Next() {
	// 	err = rows.Scan(&vv)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }
	// if err := pgxscan.ScanOne(&lol, rows); err != nil {
	// 	// Handle rows processing error.
	// }
	// mdr := lol.Elem().Interface()
	products, err := pgx.CollectRows(rows, pgx.RowToStructByName[entityType])
	if err != nil {
		fmt.Printf("CollectRows error: %v", err)
		return nil, err
	}
	fmt.Printf("%v\n", lol)
	// fmt.Printf("%v\n", mdr)
	return values, nil
}
