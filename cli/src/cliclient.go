package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
)

type Reader interface {
	ReadString(delim byte) (string, error)
}

type Network interface {
	Dial(network, address string) (net.Conn, error)
}

type FileReader interface {
	ReadFile(filename string) ([]byte, error)
}

func main() {
	var fileReader FileReader = &RealFileReader{}
	var network Network = &RealNetwork{}
	var ioReader Reader = bufio.NewReader(os.Stdin)
	var client Client = &RealClient{}
	cliClient(&client, &ioReader, &network, &fileReader)
}

type Client interface {
	SetupConnection(reader *Reader, network *Network) (net.Conn, error)
	MakeConnectionReader(conn *net.Conn) *Reader
	ListenToServer(reader *Reader)
	GetMessageFromUser(reader *Reader) (string, error)
	HandleMessage(message string, fileReader *FileReader) (string, error)
	SendMessage(conn *io.Writer, rpc string)
	ConnectionValid(conn *io.Writer, reader *Reader) bool
}

func cliClient(
	client *Client,
	ioReader *Reader,
	network *Network,
	fileReader *FileReader,
) {
	for {
		initShell()
		conn, err := (*client).SetupConnection(ioReader, network)
		if err != nil {
			printError(err)
			continue
		}
		connReader := (*client).MakeConnectionReader(&conn)
		var writer io.Writer = conn
		if !(*client).ConnectionValid(&writer, connReader) {
			printError(errors.New("No node listening at that address"))
			continue
		}
		go (*client).ListenToServer(connReader)
		for {
			message, err := (*client).GetMessageFromUser(ioReader)
			if err != nil {
				printError(err)
				break
			}
			rpc, err := (*client).HandleMessage(message, fileReader)
			if err != nil {
				printError(err)
				break
			}
			fmt.Println("Sending message: ", rpc)
			(*client).SendMessage(&writer, rpc)
			if rpc == "close;" || rpc == "exit;"{
				break
			}
		}
	}
}

func printError(err error) {
	fmt.Println("\nError:\n", err.Error())
}
