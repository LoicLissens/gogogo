// FORKED FROM : https://github.com/Nexidian/gocliselectc
package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/buger/goterm"
	"github.com/pkg/term"
)

// Raw input keycodes
// Escape code list : https://www.climagic.org/mirrors/VT100_Escape_Codes.html

var (
	up     byte = 65
	down   byte = 66
	escape byte = 27
	enter  byte = 13
)
var keys = map[byte]bool{
	up:   true,
	down: true,
}

type Menu struct {
	Prompt      string
	CursorPos   int
	MenuItems   []*MenuItem
	clearScreen bool
}

type MenuItem struct {
	Text     string
	ID       string
	SubMenu  *Menu
	Disabled bool
	Color    int // goterm color
}

func NewMenu(prompt string, clearScreen bool) *Menu {
	return &Menu{
		Prompt:      prompt,
		MenuItems:   make([]*MenuItem, 0),
		clearScreen: clearScreen,
	}
}

// AddItem will add a new menu option to the menu list
func (m *Menu) AddItem(option string, id string) *Menu {
	menuItem := &MenuItem{
		Text:     option,
		ID:       id,
		Disabled: false,
		Color:    goterm.YELLOW,
	}

	m.MenuItems = append(m.MenuItems, menuItem)
	return m
}
func (m *Menu) AddCustomItem(menuItem *MenuItem) *Menu {
	m.MenuItems = append(m.MenuItems, menuItem)
	return m
}

// renderMenuItems prints the menu item list.
// Setting redraw to true will re-render the options list with updated current selection.
func (m *Menu) renderMenuItems(redraw bool) {
	if redraw {
		// Move the cursor up n lines where n is the number of options, setting the new
		// location to start printing from, effectively redrawing the option list
		//
		// This is done by sending a VT100 escape code to the terminal
		// @see http://www.climagic.org/mirrors/VT100_Escape_Codes.html
		fmt.Printf("\033[%dA", len(m.MenuItems)-1)
	}

	for index, menuItem := range m.MenuItems {
		var newline = "\n"
		if index == len(m.MenuItems)-1 {
			// Adding a new line on the last option will move the cursor position out of range
			// For out redrawing
			newline = ""
		}
		menuItemText := menuItem.Text
		var color int
		if menuItem.Disabled {
			color = goterm.WHITE
		} else {
			color = goterm.YELLOW
		}
		cursor := "  "
		if index == m.CursorPos {
			cursor = goterm.Color("> ", goterm.YELLOW)
			menuItemText = goterm.Color(menuItemText, color)
		}

		fmt.Printf("\r%s %s%s", cursor, menuItemText, newline)
	}
}

// Display will display the current menu options and awaits user selection
// It returns the users selected choice
func (m *Menu) Display() string {

	// Show cursor again.
	defer showCursor()
	if m.clearScreen {
		clearScreen()
	}

	fmt.Printf("%s\n", goterm.Color(goterm.Bold(m.Prompt)+":", goterm.CYAN))

	m.renderMenuItems(false)

	// Turn the terminal cursor off
	fmt.Printf("\033[?25l")

	for { //awaiting for user interaction
		keyCode := getInput()
		if keyCode == escape {
			fmt.Printf("\n" + goterm.Color(goterm.Bold("Exit the application..."), goterm.GREEN))
			showCursor() // need to show cursor again before exiting as os.Exit will not reset the terminal state and execute defer functions
			os.Exit(0)
		} else if keyCode == enter {
			menuItem := m.MenuItems[m.CursorPos]
			if !menuItem.Disabled {
				fmt.Println("\r")
				return menuItem.ID
			}
		} else if keyCode == up {
			m.CursorPos = (m.CursorPos + len(m.MenuItems) - 1) % len(m.MenuItems)
			m.renderMenuItems(true)
		} else if keyCode == down {
			m.CursorPos = (m.CursorPos + 1) % len(m.MenuItems)
			m.renderMenuItems(true)
		}
	}
}

// getInput will read raw input from the terminal
// It returns the raw ASCII value inputted
func getInput() byte {
	t, _ := term.Open("/dev/tty") // Open a tty session in raw mode

	err := term.RawMode(t)
	if err != nil {
		log.Fatal(err)
	}

	var read int
	readBytes := make([]byte, 3)
	read, err = t.Read(readBytes)

	t.Restore()
	t.Close()

	// Arrow keys are prefixed with the ANSI escape code which take up the first two bytes.
	// The third byte is the key specific value we are looking for.
	// For example the left arrow key is '<esc>[A' while the right is '<esc>[C'
	// See: https://en.wikipedia.org/wiki/ANSI_escape_code
	// A ctrl key (up,down,lef,right) is 3 bytes length and the specific key is in the third byte.
	if read == 3 {
		if _, ok := keys[readBytes[2]]; ok {
			return readBytes[2]
		}
	} else {
		return readBytes[0]
	}

	return 0
}
func clearScreen() {
	// \033[H : Move the cursor to the top left corner
	// \033[2J : Clear the screen
	fmt.Print("\033[H\033[2J")
}
func showCursor() {
	fmt.Printf("\033[?25h")
}
