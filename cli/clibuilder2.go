// TODO Bug inside
package cli

// import (
// 	"fmt"
// 	"log"

// 	"github.com/pkg/term"
// )

// // Escape code list : https://www.climagic.org/mirrors/VT100_Escape_Codes.html
// var (
// 	up     byte = 65
// 	down   byte = 66
// 	escape byte = 27
// 	enter  byte = 13
// )

// type MenuItem struct {
// 	Text string
// 	ID   string
// }
// type Menu struct {
// 	Prompt         string
// 	CursorPosition int
// 	MenuItems      []*MenuItem
// }

// func get_format_input(up byte, down byte) byte {
// 	var keys = map[byte]bool{
// 		up:   true,
// 		down: true,
// 	}
// 	t, _ := term.Open("/dev/tty") // Open a tty session in raw mode
// 	err := term.RawMode(t)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	var read int
// 	bytesArray := make([]byte, 3)
// 	read, err = t.Read(bytesArray)

// 	t.Restore()
// 	t.Close()
// 	if read == 3 {
// 		if _, isKey := keys[bytesArray[2]]; isKey { // A ctrl key (up,down,lef,right) is 3 bytes length and the specific key is in the third byte.
// 			return bytesArray[2]
// 		}

// 	} else {
// 		return bytesArray[0] // other key input
// 	}
// 	return 0
// }

// func NewMenu(prompt string) *Menu {
// 	return &Menu{
// 		Prompt:    prompt,
// 		MenuItems: make([]*MenuItem, 0),
// 	}
// }
// func (self *Menu) AddItem(option string, id string) *Menu {
// 	menuItem := &MenuItem{
// 		Text: option,
// 		ID:   id,
// 	}
// 	self.MenuItems = append(self.MenuItems, menuItem)
// 	return self
// }
// func (m *Menu) renderMenuItems(redraw bool) {
// 	if redraw {
// 		// Move the cursor up n lines where n is the number of options, setting the new
// 		// location to start printing from, effectively redrawing the option list
// 		//
// 		// This is done by sending a VT100 escape code to the terminal
// 		// @see http://www.climagic.org/mirrors/VT100_Escape_Codes.html
// 		fmt.Printf("\033[%dA", len(m.MenuItems)-1)
// 	}

// 	for index, menuItem := range m.MenuItems {
// 		var newline = "\n"
// 		if index == len(m.MenuItems)-1 {
// 			// Adding a new line on the last option will move the cursor position out of range
// 			// For out redrawing
// 			newline = ""
// 		}
// 		var menuItemText string
// 		cursor := "  "
// 		if index == m.CursorPosition {
// 			cursor = "> "
// 			menuItemText = menuItem.Text
// 		}

// 		fmt.Printf("\r%s %s%s", cursor, menuItemText, newline)
// 	}
// }
// func (m *Menu) Draw_prompt() string {
// 	defer func() {
// 		// Show cursor again.
// 		fmt.Printf("\033[?25h")
// 	}()

// 	fmt.Printf("%s\n", m.Prompt+":")

// 	m.renderMenuItems(false)

// 	for { //awaiting for user interaction

// 		switch input := get_format_input(up, down); input {
// 		case up:
// 			m.CursorPosition = (m.CursorPosition + len(m.MenuItems) - 1) % len(m.MenuItems)
// 			m.renderMenuItems(true)
// 		case down:
// 			m.CursorPosition = (m.CursorPosition + 1) % len(m.MenuItems)
// 			m.renderMenuItems(true)
// 		case escape:
// 			println("Exit the application...")
// 			return ""
// 		case enter:
// 			menuItem := m.MenuItems[m.CursorPosition]
// 			fmt.Println("\r")
// 			return menuItem.ID
// 		}
// 	}

// }
