package cli

import (
	api_p "api"
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	hashing "utils/hashing"
	networkutils "utils/network"
	nodeutils "utils/node"
)

const (
	ConnPort = "80"
	ConnType = "tcp"
)

// Listens to its own IP on port 80
// Attempts to establish a connection to the incoming connection request.
// If a terminate('exit') command was issued, stops the listener to avoid panic errors as the node is terminated.
func CliServer(cliChannel chan string, api api_p.API) {
	connHost, err := networkutils.GetIP()
	errorPrinter(err)
	listener, err := net.Listen(ConnType, connHost+":"+ConnPort)
	listenerError(err)
	defer listener.Close()
	for {
		fmt.Println("Listening on " + connHost + ":" + ConnPort)
		conn, err := listener.Accept()
		listenerError(err)
		greetingsMessage(conn)
		s := handleRequest(conn, cliChannel, api)
		if s == "terminate" {
			return
		}
	}
}

// Simulates a command prompt and reads incoming messages from the connection.
// Splits the message into slices, a command and an argument (if there is one).
// First word is the command, rest is the argument.
func handleRequest(conn net.Conn, cliChannel chan string, api api_p.API) (st string) {
	for {
		netData, err := bufio.NewReader(conn).ReadString('\n')
		errorPrinter(err)
		slicedMessage := strings.SplitN(netData, " ", 2)
		s := commandHandler(conn, api, cliChannel, slicedMessage)
		if s == "close" {
			break
		}
		if s == "terminate" {
			return s
		}
	}
	s := ""
	return s
}

// Returns a string to check if user wants to close the connection and break the loop in "handleRequest"
// Also returns in case of terminate to avoid terminating with panic errors by the node.
// Calls on the proper function to handle the incoming command.
// slicedMessage[0] is the command, slicedMessage[1] is the argument if an argument exists.
func commandHandler(conn net.Conn, api api_p.API, cliChannel chan string, slicedMessage []string) (s string) {
	if strings.TrimSpace(slicedMessage[0]) == "close" {
		_, err := conn.Write([]byte("Closing connection.\n"))
		errorPrinter(err)
		conn.Close()
		println("Closed connection to client.")
		s := "close"
		return s
	} else if strings.TrimSpace(slicedMessage[0]) == "exit" {
		_, err := conn.Write([]byte("Terminated.\n"))
		errorPrinter(err)
		conn.Close()
		s := "terminate"
		terminate(cliChannel)
		return s
	} else if strings.TrimSpace(slicedMessage[0]) == "get" {
		getObject(conn, api, slicedMessage[1])
		// } else if strings.TrimSpace(slicedMessage[0]) == "forget" {
		// 	forgetTTL(conn, slicedMessage[1])
	} else if strings.TrimSpace(slicedMessage[0]) == "put" {
		putObject(conn, api, slicedMessage[1])
	} else if strings.TrimSpace(slicedMessage[0]) == "ping" {
		ping(conn, api, slicedMessage[1])
	} else if strings.TrimSpace(slicedMessage[0]) == "help" {
		availableCommands(conn)
	} else {
		_, err := conn.Write([]byte("Invalid command, for a list of commands enter 'help'.\n"))
		errorPrinter(err)
	}
	return
}

// Greeting message to new client
func greetingsMessage(conn net.Conn) {
	_, err := conn.Write([]byte("For a list of supported commands, please enter 'help'.\n"))
	errorPrinter(err)
}

// Simply lists all supported commands and how to use them.
func availableCommands(conn net.Conn) {
	_, err := conn.Write([]byte("List of supported commands:," +
		"'close' Closes the connection to this node.," +
		"'exit' Terminates the node.," +
		"'get hashNr' hashNr is an argument and returns its stored value if one exists.," +
		"'forget hashNr' NOT YET IMPLEMENTED.," +
		"'ping ipAddr' instruct the node to try and ping the given ipAddr.," +
		"'put filename' reads the contents of the given file and attempts to store it on a kademlia node.," +
		"Example: put test.txt\n"))
	errorPrinter(err)
}

// Terminates the node by sending "exit" to the loop in main.go
func terminate(cliChannel chan string) {
	fmt.Println("Terminating...")
	cliChannel <- "exit"
}

// Should return its hash if successfully restored, currently only returns wether the store was successful or not.
func putObject(conn net.Conn, api api_p.API, value string) {
	key := api.Store(value)
    _, err := conn.Write([]byte("Stored at: " + key.String()))
	errorPrinter(err)
}

func getObject(conn net.Conn, api api_p.API, hashNr string) {
	fmt.Println("Retreiving...")
    value, err := api.FindValue(hashing.ToKademliaID(hashNr))
    if err != nil {
        errorPrinter(err)
        return
    }
    _, err = conn.Write([]byte("Value: " + value))
	errorPrinter(err)
}

func ping(conn net.Conn, api api_p.API, ipAddr string) {
	fmt.Println("Pinging...")
	node := nodeutils.Node{IP: ipAddr}
    ok := api.Ping(&node)
    _, err := conn.Write([]byte("Online: " + strconv.FormatBool(ok) + "\n"))
	errorPrinter(err)
}

// NOT YET IMPLEMENTED FUNCTIONALITY
// func forgetTTL(conn net.Conn, hashNr string) {
// 	fmt.Println("Stopping refresh...")
// 	_, err := conn.Write([]byte("You want to stop refreshing: " + hashNr + "\n"))
// 	errorPrinter(err)
// }

// ERROR HELPERS
func errorPrinter(err error) {
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
}

func listenerError(err error) {
	if err != nil {
		fmt.Println("Error accepting: ", err.Error())
		os.Exit(1)
	}
}
