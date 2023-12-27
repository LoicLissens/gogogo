package db

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GuildeTable struct {
	Uuid       uuid.UUID `db:"uuid" sql_properties:"UUID PRIMARY KEY"`
	Name       string    `db:"name" sql_properties:"VARCHAR(255) NOT NULL"`
	Img_url    string    `db:"img_url" sql_properties:"VARCHAR(255)"`
	Page_url   string    `db:"page_url" sql_properties:"VARCHAR(255)"`
	Created_at time.Time `db:"created_at" sql_properties:"TIMESTAMP NOT NULL"`
	Updated_at time.Time `db:"updated_at" sql_properties:"TIMESTAMP NOT NULL"`
}

func (g GuildeTable) GetTag(fieldName string, tagName string) string {
	t := reflect.TypeOf(g)
	field, found := t.FieldByName(fieldName)
	if !found {
		return ""
		//TODO: Error
	}
	tag := field.Tag.Get(tagName)
	return tag
}

func (g GuildeTable) getDBColumnName(fieldName string) string {
	return g.GetTag(fieldName, "db")
}
func (g GuildeTable) getDBColumnProperties(fieldName string) string {
	return g.GetTag(fieldName, "sql_properties")
}
func (g GuildeTable) getTableName() string {
	return "guildes"
}
func (g GuildeTable) getDBColumn(fieldName string) string {
	name := g.getDBColumnName(fieldName)
	properties := g.getDBColumnProperties(fieldName)
	return fmt.Sprintf("%s %s,", name, properties)
}

func (g GuildeTable) CreateTable(conn *pgxpool.Pool) {
	statement := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s", g.getTableName())
	typ := reflect.TypeOf(g)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		statement += fmt.Sprintf(", %s", g.getDBColumn(field.Name))
		if i == typ.NumField() {
			statement += ";"
		}
	}
	conn.Exec(context.Background(), statement)
	log.Printf("%s table was added.", g.getTableName())
}
