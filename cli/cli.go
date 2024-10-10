package cli

import (
	"fmt"
	"jiva-guildes/backend"
	"jiva-guildes/domain/commands"
	"jiva-guildes/domain/ports/views/dtos"
	"reflect"
	"time"

	"github.com/buger/goterm"
	"github.com/google/uuid"
)

const (
	BROWSE_GUILDE = "BROWSE_GUILDE"
	CREATE_GUILDE = "CREATE_GUILDE"
	NEXT_PAGE     = "NEXT_PAGE"
	PREV_PAGE     = "PREV_PAGE"
	QUIT          = "QUIT"
)

func Manage() {
	menu := NewMenu("Manage", true)
	menu.AddItem("browse guilde", BROWSE_GUILDE)
	menu.AddItem("create guilde", CREATE_GUILDE)
	itemId := menu.Display()
	switch itemId {
	case BROWSE_GUILDE:
		browseGuildes(1)
	case CREATE_GUILDE:
		createGuilde()
	}
}
func browseGuildes(page int) {
	viewManager := backend.ViewsManager
	limit := 15
	guildesListDTO, err := viewManager.Guilde().List(page, limit)
	if err != nil {
		panic(err)
	}
	items, nbItems := guildesListDTO.Items, guildesListDTO.NbItems
	menu := NewMenu("Total guilde: "+fmt.Sprint(nbItems), true)
	menu.AddItem(goterm.Color("⬅︎ Previous page", goterm.MAGENTA), PREV_PAGE)
	for _, guilde := range items {
		menu.AddItem(guilde.Name, guilde.Uuid.String())
	}
	menu.AddItem(goterm.Color("Next page ⮕", goterm.MAGENTA), NEXT_PAGE)
	menu.AddItem(goterm.Color("Quit", goterm.RED), QUIT)
	itemId := menu.Display()
	switch itemId {
	case NEXT_PAGE:
		browseGuildes(page + 1)
	case PREV_PAGE:
		browseGuildes(page - 1)
	case QUIT:
		return
	default:
		manageGuilde(findGuilde(items, itemId))
	}
	browseGuildes(page)
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
	itemId := menu.Display()
	switch itemId {
	case "EDIT":
		editGuilde(guilde)
	case "DELETE":
		deleteGuilde(guilde.Uuid)
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
