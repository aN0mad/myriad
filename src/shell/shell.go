package shell

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

type Shell struct {
	Prompt string
	Conn   net.Conn
}

// Shell "constructor"
func StartShell(inputConn net.Conn) {
	shell := new(Shell)
	shell.Prompt = "Shell> "
	shell.Conn = inputConn

	shell.HandleShell()
}

// Function to check if channel is closed
func IsClosed(ch <-chan int) bool {
	select {
	case <-ch:
		return true
	default:
	}

	return false
}

// Send command to conn
func (shell Shell) sendCommand(cmd string) {
	cmd = strings.Replace(cmd, "\r\n", "\n", -1)
	shell.Conn.Write([]byte(cmd))
}

// Read return data from the shell with timeout
func (shell Shell) readReturn(timeout int) {
	conn := shell.Conn
	reader := bufio.NewReader(conn)
	for {
		// set SetReadDeadline
		err := conn.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
		if err != nil {
			log.Println("SetReadDeadline failed:", err)
			return
		}

		data, _, err := reader.ReadLine()
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				// time out
				return
			} else {
				fmt.Println("[!] Error reading:", err.Error())
			}
		}
		fmt.Println(string(data))
	}
}

// Main shell loop to handle prompt and interaction
func (shell Shell) HandleShell() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("\n%s", shell.Prompt)
		userInput, err := reader.ReadString('\n')
		userInput = strings.Replace(userInput, "\r\n", "", -1)
		userInput = strings.Replace(userInput, "\r", "", -1)
		if err != nil {
			fmt.Println("[!] Error input:", err.Error())
		}
		if userInput == "\r\n" || userInput == "\n" || userInput == "" {
			continue
		}
		if userInput == "background" {
			fmt.Println("Leaving shell now")
			return
		}
		shell.sendCommand(userInput + "\n")
		shell.readReturn(2)
	}
}
