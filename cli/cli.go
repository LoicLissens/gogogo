package cli

import (
	"fmt"
	"jiva-guildes/backend"
	"jiva-guildes/domain/commands"
	"jiva-guildes/domain/ports/views"
	"jiva-guildes/domain/ports/views/dtos"
	"os"
	"reflect"
	"time"

	"github.com/buger/goterm"
	"github.com/google/uuid"
)

const (
	BROWSE_GUILDES    = "BROWSE_GUILDES"
	BROWSE_ALL        = "BROWSE_ALL"
	BROWSE_VERIFIED   = "BROWSE_VERIFIED"
	BROWSE_UNVERIFIED = "BROWSE_UNVERIFIED"
	CREATE_GUILDE     = "CREATE_GUILDE"
	NEXT_PAGE         = "NEXT_PAGE"
	PREV_PAGE         = "PREV_PAGE"
	RETURN            = "RETURN"
	QUIT              = "QUIT"
)

func Manage() {
	menu := NewMenu("Manage", true)
	menu.AddItem("Browse guilde", BROWSE_GUILDES)
	menu.AddItem("Create guilde", CREATE_GUILDE)
	itemId := menu.Display()
	switch itemId {
	case BROWSE_GUILDES:
		browseGuildes(views.ListGuildesViewOpts{})
	case CREATE_GUILDE:
		createGuilde()
	}
}
func browseGuildes(filters views.ListGuildesViewOpts) {
	viewManager := backend.ViewsManager
	if filters.Limit == 0 { // means no filters
		filters = views.ListGuildesViewOpts{Page: 1, Limit: 15}
		menu := NewMenu("What would you retrieve", true)
		menu.AddItem("All guildes", BROWSE_ALL)
		menu.AddItem("Verified guildes", BROWSE_VERIFIED)
		menu.AddItem("Unverified guildes", BROWSE_UNVERIFIED)
		itemId := menu.Display()
		switch itemId {
		case BROWSE_ALL:
			browseGuildes(filters) // TODO: Do not call the same function over and over to avoid creating too much stack
		case BROWSE_VERIFIED:
			validated := true
			filters.Validated = &validated
			browseGuildes(filters)
		case BROWSE_UNVERIFIED:
			validated := false
			filters.Validated = &validated
			browseGuildes(filters)
		}

	}
	guildesListDTO, err := viewManager.Guilde().List(filters)
	if err != nil {
		panic(err)
	}
	items, nbItems := guildesListDTO.Items, guildesListDTO.NbItems
	totalPages := nbItems / filters.Limit
	currentPage := filters.Page
	menu := NewMenu(fmt.Sprintf("Total guilde: %d (page %d/%d)", nbItems, currentPage, totalPages), true)
	menu.AddCustomItem(&MenuItem{
		Text:     "⬅︎ Previous page",
		ID:       PREV_PAGE,
		Disabled: currentPage == 1,
		Color:    goterm.MAGENTA,
	})
	for _, guilde := range items {
		menu.AddItem(guilde.Name, guilde.Uuid.String())
	}
	menu.AddCustomItem(&MenuItem{
		Text:     "Next page ⮕",
		ID:       NEXT_PAGE,
		Disabled: currentPage == totalPages,
		Color:    goterm.MAGENTA,
	})
	menu.AddItem(goterm.Color("Quit", goterm.RED), QUIT)
	itemId := menu.Display()
	switch itemId {
	case NEXT_PAGE:
		filters.Page += 1
		browseGuildes(filters)
	case PREV_PAGE:
		filters.Page -= 1
		browseGuildes(filters)
	case QUIT:
		os.Exit(0) //TODO: Shutdown gracefully
	default:
		manageGuilde(findGuilde(items, itemId))
	}
	browseGuildes(filters)
}

func createGuilde() { // Berk
	fmt.Println(goterm.Color(goterm.Bold("Creating guilde"), goterm.CYAN))
	cmd := commands.CreateGuildeCommand{}

	val := reflect.ValueOf(&cmd).Elem()
	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		fmt.Println(goterm.Color(fmt.Sprintf("%s\n :", typeField.Name), goterm.YELLOW))
		var input string
		switch typeField.Type.Kind() {
		case reflect.String:
			fmt.Scanln(&input)
			val.Field(i).SetString(input)
		case reflect.Bool:
			fmt.Println(goterm.Color("\n(Enter true, false)", goterm.YELLOW))
			fmt.Scanln(&input)
			val.Field(i).SetBool(*parseBoolFromInput(input))
		case reflect.Pointer:
			// Check if it's a pointer to a bool
			if typeField.Type.Elem().Kind() == reflect.Bool {
				fmt.Println(goterm.Color("\n(Enter true, false or pass)", goterm.YELLOW))
				fmt.Scanln(&input)
				val.Field(i).Set(reflect.ValueOf(parseBoolFromInput(input)))
			}
			// Check if it's a pointer to a time.Time
			if typeField.Type.Elem() == reflect.TypeOf(time.Time{}) {
				fmt.Println(goterm.Color("Enter date (YYYY-MM-DD)", goterm.YELLOW))
				fmt.Scanln(&input)
				if input != "" {
					timeInput, err := time.Parse("2006-01-02", input)
					if err != nil {
						panic(err)
					}
					val.Field(i).Set(reflect.ValueOf(&timeInput))
				}
			}
		}
	}
	if err := backend.Validate.Struct(cmd); err != nil {
		panic(err)
	}
	g, err := backend.ServiceManager.CreateGuildeHandler(cmd)
	if err != nil {
		panic(err)
	}
	fmt.Println(goterm.Color("Guilde created: "+g.Name, goterm.GREEN))
	Manage()
}

func parseBoolFromInput(input string) *bool {
	var result bool
	switch input {
	case "true":
		result = true
		return &result
	case "false":
		result = false
		return &result
	default:
		return nil
	}
}

func manageGuilde(guilde dtos.GuildeViewDTO) {
	menu := NewMenu("Manage guilde: "+fmt.Sprint(guilde.Name), true)
	menu.AddItem("Edit", "EDIT")
	menu.AddItem("Delete", "DELETE")
	menu.AddItem("Return", RETURN)
	itemId := menu.Display()
	switch itemId {
	case "EDIT":
		editGuilde(guilde)
	case "DELETE":
		deleteGuilde(guilde.Uuid)
	case RETURN:
		return
	}

}

func findGuilde(slice []dtos.GuildeViewDTO, uuidString string) dtos.GuildeViewDTO {
	uuid, _ := uuid.Parse(uuidString)
	for i, item := range slice {
		uuidItem := item.Uuid
		if uuidItem == uuid {
			return slice[i]
		}
	}
	panic("Not found")
}

func deleteGuilde(uuid uuid.UUID) {
	cmd := commands.DeleteGuildeCommand{Uuid: uuid}
	if err := backend.Validate.Struct(cmd); err != nil {
		panic(err)
	}
	err := backend.ServiceManager.DeleteGuildeHandler(cmd)
	if err != nil {
		panic(err)
	}
	deletedMsg := goterm.Color("Guilde deleted", goterm.GREEN)
	fmt.Println(deletedMsg)
}

func editGuilde(guilde dtos.GuildeViewDTO) {
	vals := reflect.ValueOf(guilde)
	cmd := commands.UpdateGuildeCommand{Uuid: guilde.Uuid}
	menu := NewMenu("Edit guilde: "+fmt.Sprint(guilde.Name), true)

	type MappingType struct {
		Name string
		Func func(string)
	}
	editMapping := map[string]MappingType{
		"Name":     {Name: "Name", Func: func(input string) { cmd.Name = input }},
		"Img_url":  {Name: "Img_url", Func: func(input string) { cmd.Img_url = input }},
		"Page_url": {Name: "Page_url", Func: func(input string) { cmd.Page_url = input }},
		"Exists":   {Name: "Exists", Func: func(input string) { cmd.Exists = parseBoolFromInput(input) }},
		"Active":   {Name: "Active", Func: func(input string) { cmd.Active = parseBoolFromInput(input) }},
		"Creation_date": {Name: "Creation_date", Func: func(input string) {
			timeInput, err := time.Parse("2006-01-02", input)
			if err != nil {
				panic(err)
			}
			cmd.CreationDate = timeInput
		}},
		"Validated": {Name: "Validated", Func: func(input string) { cmd.Validated = parseBoolFromInput(input) }},
	}

	for i := 0; i < vals.NumField(); i++ {
		field := vals.Type().Field(i)
		value := vals.Field(i).Interface()
		_, ok := editMapping[field.Name]
		if !ok {
			prompt := field.Name + ": " + fmt.Sprint(value)
			menu.AddCustomItem(&MenuItem{
				Text:     prompt,
				ID:       field.Name,
				Disabled: true,
				Color:    goterm.WHITE,
			})
		} else {
			prompt := goterm.Bold(field.Name+": ") + fmt.Sprint(value)
			menu.AddItem(prompt, field.Name)
		}
	}
	itemId := menu.Display()
	fmt.Println(goterm.Color("Enter new value for the field : "+itemId, goterm.YELLOW))
	var input string
	fmt.Scanln(&input)
	editMapping[itemId].Func(input)
	if err := backend.Validate.Struct(cmd); err != nil {
		panic(err)
	}
	_, err := backend.ServiceManager.UpdateGuildeHandler(cmd)
	if err != nil {
		panic(err)
	}
	fmt.Println(goterm.Color("Guilde Updted", goterm.GREEN))
}
