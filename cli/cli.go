package cli

import (
	"fmt"
	"jiva-guildes/backend"
	"jiva-guildes/domain/commands"
	"jiva-guildes/domain/ports/views/dtos"
	"reflect"

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
	menu := NewMenu("Manage")
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
	menu := NewMenu("Total guilde: " + fmt.Sprint(nbItems))
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
func createGuilde() {
	fmt.Println(goterm.Color(goterm.Bold("Creating guilde"), goterm.CYAN))
	cmd := commands.CreateGuildeCommand{}

	val := reflect.ValueOf(&cmd).Elem()
	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		fmt.Printf(goterm.Color(fmt.Sprintf("%s\n :", typeField.Name), goterm.YELLOW))
		var input string
		switch typeField.Type.Kind() {
		case reflect.String:
			fmt.Scanln(&input)
			val.Field(i).SetString(input)
		case reflect.Bool:
			val.Field(i).SetBool(input)
		}
	}
	// if err := backend.Validate.Struct(cmd); err != nil {
	// 	panic(err)
	// }
}
func manageGuilde(guilde dtos.GuildeViewDTO) {
	menu := NewMenu("Manage guilde: " + fmt.Sprint(guilde.Name))
	menu.AddItem("Edit", "EDIT")
	menu.AddItem("Delete", "DELETE")
	itemId := menu.Display()
	switch itemId {
	case "EDIT":
		editGuilde(guilde.Uuid)
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
func editGuilde(uuid uuid.UUID) {
	fmt.Println("Editing guilde")
}
