package main

import (
	"flag"
	tables "jiva-guildes/adapters/db/tables"
	"jiva-guildes/backend/router"
	"jiva-guildes/backend/scripts"
	"jiva-guildes/cli"
	"jiva-guildes/scrapper"
)

type Actions int
type ActionFunction func()

const (
	SCRAP Actions = iota
	SERVE
	INIT_DB
	POPULATE_FROM_CSV
	MANAGE
)

func (action Actions) ActionsEnum() string {
	return []string{"SCRAP", "SERVE", "INIT_DB", "POPULATE_FROM_CSV", "MANAGE"}[action]
}

func main() {
	actionMapper := map[string]ActionFunction{
		SCRAP.ActionsEnum():             scrapper.Scrap,
		INIT_DB.ActionsEnum():           tables.InitAllTables,
		SERVE.ActionsEnum():             router.Serve,
		POPULATE_FROM_CSV.ActionsEnum(): scripts.PopulateDBFromCSV,
		MANAGE.ActionsEnum():            cli.Manage,
	}
	isCliMode := flag.Bool("cli", false, "Wether the module should be launched in CLI mode.")
	flag.Parse()

	if *isCliMode {
		menu := cli.NewMenu("What do you want to do ?", true)
		menu.AddItem("Manage", MANAGE.ActionsEnum())
		menu.AddItem("Init database", INIT_DB.ActionsEnum())
		menu.AddItem("Scrapping of data", SCRAP.ActionsEnum())
		menu.AddItem("Serve", SERVE.ActionsEnum())
		menu.AddItem("Populate from CSV", POPULATE_FROM_CSV.ActionsEnum())

		itemId := menu.Display()
		action := actionMapper[itemId]
		action()
	} else {
		router.Serve()
	}
}
