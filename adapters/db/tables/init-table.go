package tables

import (
	"jiva-guildes/adapters/db"
)

var allTables = []Table{&GuildeTable{}}

func InitAllTables() {
	connectionPool := db.MountDB()
	for _, val := range allTables {
		val.CreateTable(connectionPool)
	}
	db.Teardown(connectionPool)
}
