package main

import (
	"flag"
	"jiva-guildes/backend/router"
	"jiva-guildes/cli"
	tables "jiva-guildes/db/tables"
	"jiva-guildes/scrapper"
)

type Actions int
type ActionFunction func()

const (
	SCRAP Actions = iota
	SERVE
	INIT_DB
)

func (action Actions) ActionsEnum() string {
	return []string{"SCRAP", "SERVE", "INIT_DB"}[action]
}

func main() {
	actionMapper := map[string]ActionFunction{
		SCRAP.ActionsEnum():   scrapper.Scrap,
		INIT_DB.ActionsEnum(): tables.InitAllTables,
		SERVE.ActionsEnum():   router.Serve,
	}
	isCliMode := flag.Bool("cli", false, "Wether the module should be launched in CLI mode.")
	flag.Parse()

	if *isCliMode == true {
		menu := cli.NewMenu("What do you want to do ?")
		menu.AddItem("Scrapping of data", SCRAP.ActionsEnum())
		menu.AddItem("Init database", INIT_DB.ActionsEnum())
		menu.AddItem("Serve", SERVE.ActionsEnum())
		itemId := menu.Draw_prompt()
		action := actionMapper[itemId]
		action()
	}
}
