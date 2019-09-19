package cli

import (
	"api"
	"bytes"
	"fmt"
	"net"
	"os"
	"utils/hashing"
	networkutils "utils/network"
)

const (
	ConnPort = "80"
	ConnType = "tcp"
)

func CliApp() {
	// get IP for the listener.
	connHost, err := networkutils.GetIP()
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	// Listen for incoming connections.
	listener, err := net.Listen(ConnType, connHost+":"+ConnPort)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer listener.Close()
	fmt.Println("Listening on " + connHost + ":" + ConnPort)
	for {
		// Listen for an incoming connection.
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	// Translate bytes to string.
	n := bytes.Index(buf, []byte{0})
	message := string(buf[:n-1])
	if message == "exit" {
		terminate()
		_, err := conn.Write([]byte("Terminated.\n"))
		if err != nil {
			fmt.Println("Error reading:", err.Error())
		}
	} else if message == "put" {
		putObject()
		_, err := conn.Write([]byte("Uploaded.\n"))
		if err != nil {
			fmt.Println("Error reading:", err.Error())
		}
	} else if message == "get" {
		getObject()
		_, err := conn.Write([]byte("Retrieved.\n"))
		if err != nil {
			fmt.Println("Error reading:", err.Error())
		}
	} else {
		_, err := conn.Write([]byte("Please enter the command 'exit', 'put' or 'get'.\n"))
		if err != nil {
			fmt.Println("Error reading:", err.Error())
		}
	}
	err = conn.Close()
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
}

func terminate() {
	fmt.Println("Terminating...")
}

func putObject() {
	api.Store("hello")
}

func getObject() {
	api.FindValue(hashing.NewRandomKademliaID())
}
