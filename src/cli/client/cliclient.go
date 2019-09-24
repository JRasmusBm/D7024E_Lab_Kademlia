package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
)


func main() {
	cliClient()
}

// Starts the shell and attempts to make a connection to the given ipAddr:port
// Calls functions to handle sending and receiving messages.
func cliClient() {
	reader := bufio.NewReader(os.Stdin)
	initShell()
	nodeConnection, _ := reader.ReadString('\n')
	nodeConnection = strings.Replace(nodeConnection, "\n", "", -1)
	conn, err := net.Dial("tcp", nodeConnection)
	errorHelper(err)
	fmt.Println("Successfully connected to: " + nodeConnection)
	messageReceiver(conn, "")
	for {
		commandPrompt()
		text, _ := reader.ReadString('\n')
		messageSender(conn, text)
		messageReceiver(conn, text)
		slicedText := strings.SplitN(text," ",2)
		if strings.TrimSpace(slicedText[0]) == "close" {
			return
		}
		if strings.TrimSpace(slicedText[0]) == "exit" {
			return
		}
	}
}

// Fprintf sends information through the connection to the node.
// checks if it's to send the entire command or in the case of 'put' read the file contents and send.
func messageSender(conn net.Conn, text string) {
	slicedText := strings.SplitN(text, " ",2)
	if strings.TrimSpace(slicedText[0]) == "put" {
		content, err := ioutil.ReadFile(strings.TrimSpace(slicedText[1]))
		errorHelper(err)
		fmt.Fprintf(conn, slicedText[0] + " " + string(content) + "\n")
	} else {
		fmt.Fprintf(conn, text+"\n")
	}
}

// Waits for a reply from the connection.
// The buffers delimiter is '\n', so when sending several lines of information a loop is required.
// Or it will only read the first line and assume it's done.
func messageReceiver(conn net.Conn, text string) {
	slicedText := strings.SplitN(text," ",2)
	command := slicedText[0]
	if strings.TrimSpace(command) == "help" {
		message, _ := bufio.NewReader(conn).ReadString('\n')
		message = strings.Replace(message, "\n", "", -1)
		slicedMessage := strings.Split(message, ",")
		for i := 0; i<len(slicedMessage); i++ {
			fmt.Print("->: " + slicedMessage[i] + "\n")
		}
	} else {
		message, _ := bufio.NewReader(conn).ReadString('\n')
		message = strings.Replace(message, "\n", "", -1)
		fmt.Print("->: " + message + "\n")
	}
}

// PRINT HELPERS
func initShell() {
	fmt.Println("Simple Shell")
	fmt.Println("---------------------")
	fmt.Println("Please enter the ipAddr and port you wish to interact with.")
	fmt.Println("Example: 172.19.0.39:80")
	commandPrompt()
}

func commandPrompt() {
	fmt.Printf(">> ")
}

// ERROR HELPERS
func errorHelper(err error) {
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
}
