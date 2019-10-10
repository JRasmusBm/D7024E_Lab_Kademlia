package cli

import (
	api_p "api"
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	hashing "utils/hashing"
	networkutils "utils/network"
	nodeutils "utils/node"
)

type RealServer struct {
	api api_p.API
}

func (r *RealServer) SetupListener(network *Network) (net.Listener, error) {
	connHost, err := networkutils.GetIP()
	if err != nil {
		return nil, err
	}
	fmt.Println("Listening on " + connHost + ":80")
	return (*network).Listen("tcp", connHost+":80")
}

func (r *RealServer) ListenForConnection(listener net.Listener) (net.Conn, error) {
	return listener.Accept()
}

func (r *RealServer) MakeConnectionReader(conn *net.Conn) *Reader {
	var reader Reader = bufio.NewReader(*conn)
	return &reader
}

func (r *RealServer) ListenToClient(reader *Reader) string {
	message, err := (*reader).ReadString(';')
	if err != nil {
		return "close;"
	}
	return message
}

func (r *RealServer) MessageParser(incomingMessage string) []string {
	incomingMessage = strings.Replace(incomingMessage, ";", "", -1)
	parsedMessage := strings.SplitN(incomingMessage, " ", 2)
	return parsedMessage
}

func (r *RealServer) CommandHandler(parsedMessage []string) string {
	if strings.TrimSpace(parsedMessage[0]) == "close" {
		return "Closing connection.;"
	} else if strings.TrimSpace(parsedMessage[0]) == "exit" {
		return "Terminating node.;"
	} else if strings.TrimSpace(parsedMessage[0]) == "get" {
		key, err := hashing.ToKademliaID(strings.TrimSpace(parsedMessage[1]))
		if err != nil {
			return "Error: Invalid Kademlia ID;"
		}
		value, err := r.api.FindValue(key)
		if err != nil {
			return err.Error()
		}
		return "Value: " + value + ";"
	} else if strings.TrimSpace(parsedMessage[0]) == "put" {
		key, err := r.api.Store(strings.TrimSpace(parsedMessage[1]))
		if err != nil {
			return err.Error()
		}
		return "Stored at: " + key.String() + ";"
	} else if strings.TrimSpace(parsedMessage[0]) == "ping" {
		node := nodeutils.Node{IP: strings.TrimSpace(parsedMessage[1])}
		ok := r.api.Ping(&node)
		return "Online: " + strconv.FormatBool(ok) + ";"
	} else if strings.TrimSpace(parsedMessage[0]) == "connect" {
		return "ok;" + supported_commands()
	} else if strings.TrimSpace(parsedMessage[0]) == "help" {
		return supported_commands()
	} else {
		return "Invalid command, for a list of commands enter 'help'.;"
	}
}

func supported_commands() string {
	return "List of supported commands:\n" +
		"'close' Closes the connection to this node.\n" +
		"'exit' Terminates the node.\n" +
		"'get hashNr' hashNr is an argument and returns its stored value if one exists.\n" +
		"'ping ipAddr' instruct the node to try and ping the given ipAddr.\n" +
		"'put filename' reads the contents of the given file and attempts to store it on a kademlia node.\n" +
		"Example: put test.txt;"
}

func (r *RealServer) SendMessage(conn *net.Conn, cliChannel chan string, responseMessage string) {
	if responseMessage == "Closing connection.;" {
		_, err := (*conn).Write([]byte(responseMessage))
		if err != nil {
			fmt.Println("Was unable to send message.")
		}
		(*conn).Close()
		fmt.Println("Connection closed.")
	} else if responseMessage == "Terminating node.;" {
		_, err := (*conn).Write([]byte(responseMessage))
		if err != nil {
			fmt.Println("Was unable to send message.")
		}
		(*conn).Close()
		cliChannel <- "exit"
	} else {
		_, err := (*conn).Write([]byte(responseMessage))
		if err != nil {
			fmt.Println("Was unable to send message.")
		}
	}
}
