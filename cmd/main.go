package main

import (
	"bufio"
	"fmt"
	"io"
	"myriad/src/menu"
	"myriad/src/socket"
	"net"
	"os"
	"os/signal"
	"syscall"
)

// Store clients in an array
func handleClients(chanClient chan net.Conn, clientList *[]net.Conn) {
	for newClient := range chanClient {
		*clientList = append(*clientList, newClient)
	}
}

// Handler for Ctrl+C
func SetupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		//fmt.Println("\r- Ctrl+C pressed in Terminal")
		fmt.Println("Cleaning up and exiting...")
		cleanExit()
	}()
}

func cleanExit() {
	os.Exit(0)
}

func main() {
	ConnHost := "127.0.0.1"
	ConnPort := "9000"
	ConnType := "tcp"

	fmt.Println("Starting...")

	// Start the ctrl-c handler
	SetupCloseHandler()

	// List of connections
	clientList := []net.Conn{}

	// Channel to handle list of connections
	channelClients := make(chan net.Conn)
	ClientChan := channelClients

	// Thread to handle new connections
	go handleClients(channelClients, &clientList)

	// Create a socket
	fmt.Printf("Listening on %s:%s:%s\n\n", ConnType, ConnHost, ConnPort)
	go socket.CreateSocket(ConnHost, ConnPort, ConnType, ClientChan)

	// Create the menu
	userMenu := menu.NewMenu()

	// Create a new reader for stdin
	reader := bufio.NewReader(os.Stdin)

	// Main prog loop
	for {
		userMenu.PrintMenu()
		fmt.Printf("Num clients: %d\n", len(clientList))
		fmt.Printf("Menu> ")
		userInput, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("[!] Error input:", err.Error())
		}
		if userInput == "\r\n" || userInput == "\n" {
			continue
		}
		if err == io.EOF {
			cleanExit()
		}
		retValue := userMenu.HandleInput(userInput, clientList)
		if retValue == "exit" {
			cleanExit()
		}
		fmt.Printf("\n\n")
	}
}
