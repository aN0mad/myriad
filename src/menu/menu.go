package menu

import (
	"fmt"
	"myriad/src/shell"
	"net"
	"strconv"
	"strings"
)

type Menu struct {
	MenuOpts map[string]string
}

// Menu "constructor"
func NewMenu() *Menu {
	menu := new(Menu)
	menu.MenuOpts = map[string]string{
		"Help":         "Print this menu",
		"List":         "List all sessions",
		"Session {ID}": "Interact with session with id of {ID}",
		"Exit":         "Exit the handler",
	}
	return menu
}

// Print the menu options
func (menu Menu) PrintMenu() {
	fmt.Println("Menu:\n")
	for key, value := range menu.MenuOpts {
		fmt.Printf("%s\t%s\n", key, value)
	}
}

// Handle options for menu
func (menu Menu) HandleInput(userInput string, clientList []net.Conn) string {
	userWords := strings.Fields(userInput)

	switch strings.ToLower(userWords[0]) {
	case "help":

	case "list":
		fmt.Printf("\n")
		for index, item := range clientList {
			fmt.Printf("%d. %s\n", index, item.RemoteAddr().String())
		}
	case "session":
		if len(userWords) != 2 {
			fmt.Printf("\n[!] Error: There should be 2 arguments not %d", len(userWords))
			return ""
		}
		sessionID, err := strconv.Atoi(userWords[1])
		if err != nil {
			fmt.Printf("\n[!] Error: '%s' is not a valid integer", userWords[1])
			return ""
		}
		if sessionID > len(clientList) {
			fmt.Printf("\n[!] Error: '%s' is not a valid session ID", userWords[1])
			return ""
		}
		fmt.Printf("Connect user to session: %s\n", userWords[1])
		conn := clientList[sessionID]
		shell.StartShell(conn)

	case "exit":
		fmt.Printf("Exiting the program...\n")
		return "exit"
	default:
		fmt.Println("That is not an option sorry")
		return ""
	}
	return ""
}
