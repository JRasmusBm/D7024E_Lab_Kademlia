package cli

import (
	api_p "api"
	"fmt"
	"net"
	"strings"
)

type Reader interface {
	ReadString(delim byte) (string, error)
}

type Network interface {
	Listen(network, address string) (net.Listener, error)
}

type Server interface {
	SetupListener(network *Network) (net.Listener, error)
	ListenForConnection(listener net.Listener) (net.Conn, error)
	MakeConnectionReader(conn *net.Conn) *Reader
	ListenToClient(reader *Reader) string
	MessageParser(incomingMessage string) []string
	CommandHandler(parsedMessage []string) string
	SendMessage(conn *net.Conn, cliChannel chan string, responseMessage string)
}

func CliServerInit(api api_p.API, cliChannel chan string) {
	var network Network = &RealNetwork{}
	var server Server = &RealServer{api: api}
	CliServer(cliChannel, &network, &server)
}

func CliServer(
	cliChannel chan string,
	network *Network,
	server *Server) {
	listener, err := (*server).SetupListener(network)
	if err != nil {
		printError(err)
		return
	}
	//defer listener.Close() this should not have to be commented
	for {
		conn, err := (*server).ListenForConnection(listener)
		if err != nil {
			printError(err)
			return
		}
		connReader := (*server).MakeConnectionReader(&conn)
		for {
			incomingMessage := (*server).ListenToClient(connReader)
			parsedMessage := (*server).MessageParser(incomingMessage)
			go handleRequest(server, &conn, parsedMessage, cliChannel)
			if strings.TrimSpace(parsedMessage[0]) == "close" {
				break
			}
		}
	}
}

func handleRequest(
	server *Server,
	conn *net.Conn,
	parsedMessage []string,
	cliChannel chan string) {
	responseMessage := (*server).CommandHandler(parsedMessage)
	//fmt.Println("Sending message: " + responseMessage) causes mutex overflow in testing.
	(*server).SendMessage(conn, cliChannel, responseMessage)
}

func printError(err error) {
	fmt.Println("Error:", err.Error())
}
