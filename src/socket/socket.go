package socket

import (
	"fmt"
	"net"
	"os"
)

type Socket struct {
	ConnHost   string
	ConnPort   string
	ConnType   string
	ClientChan chan net.Conn
}

// Socket "constructor"
func CreateSocket(connHost string, connPort string, connType string, clientChan chan net.Conn) *Socket {
	s := new(Socket)
	s.ConnHost = connHost
	s.ConnPort = connPort
	s.ConnType = connType
	s.ClientChan = clientChan

	list, err := net.Listen(s.ConnType, s.ConnHost+":"+s.ConnPort)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	//fmt.Printf("Listening on %s:%s:%s\n", s.ConnType, s.ConnHost, s.ConnPort)
	defer list.Close()

	for {
		client, err := list.Accept()
		if err != nil {
			fmt.Println("[!] Error connecting:", err.Error())
		}

		fmt.Println("\n\nClient " + client.RemoteAddr().String() + " connected.")

		go s.handleConnection(client)
	}
}

// Send conn to channel for handling
func (s Socket) handleConnection(clientConn net.Conn) {
	s.ClientChan <- clientConn
}
