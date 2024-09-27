package tables

import (
	"jiva-guildes/adapters/db"
	"jiva-guildes/settings"
)

var allTables = []Table{&GuildeTable{}}

func InitAllTables() {
	connectionPool := db.MountDB(settings.AppSettings.DATABASE_URI)
	for _, val := range allTables {
		val.CreateTable(connectionPool)
	}
	db.Teardown(connectionPool)
}
func DropAllTables() {
	connectionPool := db.MountDB(settings.AppSettings.DATABASE_URI)
	for _, val := range allTables {
		val.DropTable(connectionPool)
	}
	db.Teardown(connectionPool)
}
