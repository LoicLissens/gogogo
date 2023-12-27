package db

import (
	"jiva-guildes/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Table interface {
	getTableName() string
	CreateTable(conn *pgxpool.Pool)
	getDBColumnName(fieldName string) string
	getDBColumnProperties(fieldName string) string
	getDBColumn(fieldName string) string
}

var allTables = []Table{&GuildeTable{}}

func InitAllTables() interface{} {
	connectionPool := db.MountDB()
	for _, val := range allTables {
		val.CreateTable(connectionPool)
	}
	db.Teardown(connectionPool)
	return nil
}
