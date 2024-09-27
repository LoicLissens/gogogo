package tables

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Table interface {
	GetTableName() string
	CreateTable(conn *pgxpool.Pool)
	DropTable(conn *pgxpool.Pool)
}

type BaseModelTable struct {
	Uuid       uuid.UUID `db:"uuid" sql_properties:"UUID PRIMARY KEY"`
	Created_at time.Time `db:"created_at" sql_properties:"TIMESTAMP NOT NULL"`
	Updated_at time.Time `db:"updated_at" sql_properties:"TIMESTAMP NOT NULL"`
}

func GetTag(fieldName string, tagName string, table Table) (string, error) {
	t := reflect.TypeOf(table)
	field, found := t.FieldByName(fieldName)
	if !found {
		return "", fmt.Errorf("field not found: %s", fieldName)
	}
	tag := field.Tag.Get(tagName)
	return tag, nil
}

func GetDBColumnName(fieldName string, table Table) string {
	tag, _ := GetTag(fieldName, "db", table)
	return tag
}
func getDBColumnProperties(fieldName string, table Table) string {
	tag, _ := GetTag(fieldName, "sql_properties", table)
	return tag
}
func (table BaseModelTable) GetTableName() (string, error) {
	return "", errors.New("GetTableName not implemented")
}

func getDBColumn(fieldName string, table Table) string {
	name := GetDBColumnName(fieldName, table)
	properties := getDBColumnProperties(fieldName, table)
	return fmt.Sprintf("%s %s", name, properties)
}

func CreateTable(conn *pgxpool.Pool, table Table) {
	tableName := table.GetTableName()
	statement := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (", tableName)
	allFields, _ := DeepFields(table)
	for i, field := range allFields {
		if i != 0 {
			statement += ","
		}
		statement += fmt.Sprintf(" %s", getDBColumn(field.Name, table))
		if i == len(allFields)-1 {
			statement += ");"
		}
	}
	_, err := conn.Exec(context.Background(), statement)
	if err != nil {
		log.Fatalf("Error while creating table %s: %v", tableName, err)

	}
	log.Printf("%s table was added.", tableName)
}
func DropTable(conn *pgxpool.Pool, table Table) {
	tableName := table.GetTableName()
	statement := fmt.Sprintf("DROP TABLE IF EXISTS %s;", tableName)
	_, err := conn.Exec(context.Background(), statement)
	if err != nil {
		log.Fatalf("Error while dropping table %s: %v", tableName, err)
	}
	log.Printf("%s table was dropped.", tableName)
}

func DeepFields(iface interface{}) ([]reflect.StructField, []reflect.Value) {
	fields := make([]reflect.StructField, 0)
	values := make([]reflect.Value, 0)
	ifv := reflect.ValueOf(iface)
	ift := reflect.TypeOf(iface)

	for i := 0; i < ift.NumField(); i++ {
		v := ifv.Field(i)
		t := ift.Field(i)
		if t.Type == reflect.TypeOf(time.Time{}) {
			fields = append(fields, t)
			values = append(values, v)
			continue
		}
		switch v.Kind() {
		case reflect.Struct:
			field, val := DeepFields(v.Interface())
			fields = append(fields, field...)
			values = append(values, val...)
		default:
			fields = append(fields, t)
			values = append(values, v)
		}
	}
	return fields, values
}
