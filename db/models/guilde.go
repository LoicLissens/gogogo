package db

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"github.com/jackc/pgx/v5/pgxpool"
)

type GuildeTable struct {
	BaseModelTable
	Name     string `db:"name" sql_properties:"VARCHAR(255) NOT NULL"`
	Img_url  string `db:"img_url" sql_properties:"VARCHAR(255)"`
	Page_url string `db:"page_url" sql_properties:"VARCHAR(255)"`
}

func (g GuildeTable) GetTag(fieldName string, tagName string) (string, error) {
	t := reflect.TypeOf(g)
	field, found := t.FieldByName(fieldName)
	if !found {
		return "", fmt.Errorf("field not found: %s", fieldName)
	}
	tag := field.Tag.Get(tagName)
	return tag, nil
}

func (g GuildeTable) getDBColumnName(fieldName string) string {
	tag, _ := g.GetTag(fieldName, "db")
	return tag
}
func (g GuildeTable) getDBColumnProperties(fieldName string) string {
	tag, _ := g.GetTag(fieldName, "sql_properties")
	return tag
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
