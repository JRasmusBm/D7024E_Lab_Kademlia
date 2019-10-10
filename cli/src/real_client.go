package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

type RealClient struct {
}

func (r *RealClient) SetupConnection(
	reader *Reader,
	network *Network,
) (net.Conn, error) {
	nodeConnection, err := (*reader).ReadString('\n')
	if err != nil {
		return nil, err
	}
	nodeConnection = strings.Replace(nodeConnection, "\n", "", -1)
	return (*network).Dial("tcp", nodeConnection)
}

func timeout(ch chan bool) {
	time.Sleep(2 * time.Second)
	ch <- false
	return
}

func read(reader *Reader, ch chan bool) {
	message, err := (*reader).ReadString(';')
	if err != nil || message != "ok;" {
		ch <- false
		return
	}
	ch <- true
	return
}

func (r *RealClient) ConnectionValid(conn *io.Writer, reader *Reader) bool {
	ch := make(chan bool, 2)
	(*conn).Write([]byte("connect;"))
	go timeout(ch)
	go read(reader, ch)
	return <-ch
}

func (r *RealClient) MakeConnectionReader(conn *net.Conn) *Reader {
	var reader Reader = bufio.NewReader(*conn)
	return &reader
}

func (r *RealClient) ListenToServer(reader *Reader) {
	for {
		message, err := (*reader).ReadString(';')
		if err != nil {
			continue
		}
		fmt.Print("->: " + message + "\n")
	}
}

func (r *RealClient) GetMessageFromUser(reader *Reader) (string, error) {
	commandPrompt()
	return (*reader).ReadString('\n')
}

func (r *RealClient) HandleMessage(message string, fileReader *FileReader) (string, error) {
	message = strings.Replace(message, "\n", "", -1)
	slicedText := strings.SplitN(message, " ", 2)
	if strings.TrimSpace(slicedText[0]) == "put" {
		content, err := (*fileReader).ReadFile(strings.TrimSpace(slicedText[1]))
		if err != nil {
			return "", err
		}
		return slicedText[0] + " " + strings.TrimSpace(string(content)) + ";", nil
	} else {
		return message + ";", nil
	}
}

type Writer interface {
	Write(str string)
}

func (r *RealClient) SendMessage(conn *io.Writer, rpc string) {
	fmt.Fprintf(*conn, rpc)
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
	time.Sleep(100 * time.Millisecond)
	fmt.Printf(">> ")
}
