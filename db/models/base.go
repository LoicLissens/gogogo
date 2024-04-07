package models

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
	getTableName() string
	CreateTable(conn *pgxpool.Pool)
	getDBColumnName(fieldName string) string
	getDBColumnProperties(fieldName string) string
	getDBColumn(fieldName string) string
}
type BaseModelTable struct {
	Uuid       uuid.UUID `db:"uuid" sql_properties:"UUID PRIMARY KEY"`
	Created_at time.Time `db:"created_at" sql_properties:"TIMESTAMP NOT NULL"`
	Updated_at time.Time `db:"updated_at" sql_properties:"TIMESTAMP NOT NULL"`
}

func (table BaseModelTable) GetTag(fieldName string, tagName string) (string, error) {
	t := reflect.TypeOf(table)
	field, found := t.FieldByName(fieldName)
	if !found {
		return "", fmt.Errorf("field not found: %s", fieldName)
	}
	tag := field.Tag.Get(tagName)
	return tag, nil
}

func (table BaseModelTable) getDBColumnName(fieldName string) string {
	tag, _ := table.GetTag(fieldName, "db")
	return tag
}

func (table BaseModelTable) getDBColumnProperties(fieldName string) string {
	tag, _ := table.GetTag(fieldName, "sql_properties")
	return tag
}

func (table BaseModelTable) getTableName() (string, error) {
	return "", errors.New("getTableName not implemented")
}
func (table BaseModelTable) getDBColumn(fieldName string) string {
	name := table.getDBColumnName(fieldName)
	properties := table.getDBColumnProperties(fieldName)
	return fmt.Sprintf("%s %s,", name, properties)
}

func (table BaseModelTable) CreateTable(conn *pgxpool.Pool) {
	tableName, err := table.getTableName()
	if err != nil {
		statement := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s", tableName)
		typ := reflect.TypeOf(table)
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			statement += fmt.Sprintf(", %s", table.getDBColumn(field.Name))
			if i == typ.NumField() {
				statement += ";"
			}
		}
		conn.Exec(context.Background(), statement)
		log.Printf("%s table was added.", tableName)
	}
}
