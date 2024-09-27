package cli

import (
	"fmt"
	"jiva-guildes/backend"
	"jiva-guildes/domain/commands"
	"jiva-guildes/domain/ports/views/dtos"

	"github.com/buger/goterm"
	"github.com/google/uuid"
)

const (
	BROWSE_GUILDE = "BROWSE_GUILDE"
	CREATE_GUILDE = "CREATE_GUILDE"
	NEXT_PAGE     = "NEXT_PAGE"
	PREV_PAGE     = "PREV_PAGE"
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
	menu.AddItem("⬅︎ Previous page", PREV_PAGE)
	for _, guilde := range items {
		menu.AddItem(guilde.Name, guilde.Uuid.String())
	}
	menu.AddItem("Next page ⮕", NEXT_PAGE)
	itemId := menu.Display()
	switch itemId {
	case NEXT_PAGE:
		browseGuildes(page + 1)
	case PREV_PAGE:
		browseGuildes(page - 1)
	default:
		manageGuilde(findGuilde(items, itemId))
	}
	browseGuildes(page)
}
func createGuilde() {
	fmt.Println("Creating guilde")
}
func manageGuilde(guilde dtos.GuildeViewDTO) {
	menu := NewMenu("Manage guilde: " + fmt.Sprint(guilde.Name))
	menu.AddItem("Edit", "EDIT")
	menu.AddItem("Delete", "DELETE")
	itemId := menu.Display()
	switch itemId {
	case "EDIT":
		fmt.Println("Not implemented")
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
