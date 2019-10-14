package cli

import (
	api_p "api"
	"errors"
	"fmt"
	"io"
	"net"
	"network"
	"testing"
	"time"
	"utils/constants"
	"utils/hashing"
	networkutils "utils/network"
	nodeutils "utils/node"
)

func TestListenToClientDoesNotCrash(t *testing.T) {
	var server Server = &RealServer{}
	var reader Reader = &MockReader{ReadStringResult: "Hello"}
	go server.ListenToClient(&reader)
	time.Sleep(300 * time.Millisecond)
}

func TestListenToClientDoesNotCrashErr(t *testing.T) {
	var server Server = &RealServer{}
	var reader Reader = &MockReader{ReadStringErr: errors.New("Random Error")}
	go server.ListenToClient(&reader)
	time.Sleep(5 * time.Millisecond)
}

// COMMANDHANDLER TESTS
// Still need get, put, ping for commandhandler. mock API...
func TestCommandHandlerClose(t *testing.T) {
	var server Server = &RealServer{}
	expected := "Closing connection.;"
	actual := server.CommandHandler([]string{"close"})
	if expected != actual {
		t.Errorf(
			"Expected %s got %s",
			expected,
			actual)
	}
}

func TestCommandHandlerExit(t *testing.T) {
	var server Server = &RealServer{}
	expected := "Terminating node.;"
	actual := server.CommandHandler([]string{"exit"})
	if expected != actual {
		t.Errorf(
			"Expected %s got %s",
			expected,
			actual)
	}
}

func TestCommandHandlerInvalid(t *testing.T) {
	var server Server = &RealServer{}
	expected := "Invalid command, for a list of commands enter 'help'.;"
	actual := server.CommandHandler([]string{" "})
	if expected != actual {
		t.Errorf(
			"Expected %s got %s",
			expected,
			actual)
	}
}

func TestCommandHandlerConnect(t *testing.T) {
	var server Server = &RealServer{}
	expected := "ok;List of supported commands:\n" +
		"'close' Closes the connection to this node.\n" +
		"'exit' Terminates the node.\n" +
		"'get hashNr' hashNr is an argument and returns its stored value if one exists.\n" +
		"'ping ipAddr' instruct the node to try and ping the given ipAddr.\n" +
		"'put filename' reads the contents of the given file and attempts to store it on a kademlia node.\n" +
		"Example: put test.txt;"
	actual := server.CommandHandler([]string{"connect"})
	if expected != actual {
		t.Errorf(
			"Expected %s got %s",
			expected,
			actual)
	}
}

func TestCommandHandlerHelp(t *testing.T) {
	var server Server = &RealServer{}
	expected := "List of supported commands:\n" +
		"'close' Closes the connection to this node.\n" +
		"'exit' Terminates the node.\n" +
		"'get hashNr' hashNr is an argument and returns its stored value if one exists.\n" +
		"'ping ipAddr' instruct the node to try and ping the given ipAddr.\n" +
		"'put filename' reads the contents of the given file and attempts to store it on a kademlia node.\n" +
		"Example: put test.txt;"
	actual := server.CommandHandler([]string{"help"})
	if expected != actual {
		t.Errorf(
			"Expected %s got %s",
			expected,
			actual)
	}
}

// MESSAGEPARSER TESTS
func TestMessageParserOneword(t *testing.T) {
	var server Server = &RealServer{}
	expected := "testword"
	actual := server.MessageParser("testword")
	if expected != actual[0] {
		t.Errorf(
			"Expected %s got %s",
			expected,
			actual)
	}
}

func TestMessageParserSeveralWords(t *testing.T) {
	var server Server = &RealServer{}
	expectedfirst := "This"
	expectedsecond := "is a test sentence."
	actual := server.MessageParser("This is a test sentence.")
	if expectedfirst != actual[0] {
		t.Errorf(
			"Expected %s got %s",
			expectedfirst,
			actual)
	}
	if expectedsecond != actual[1] {
		t.Errorf(
			"Expected %s got %s",
			expectedsecond,
			actual)
	}
}

// SETUP LISTENER TESTS

func TestSetupListenerError(t *testing.T) {
	var networkUtils networkutils.NetworkUtils = &networkutils.MockNetworkUtils{
		Err: errors.New("TestSetupListenerError"),
	}
	var server Server = &RealServer{networkUtils: &networkUtils}
	_, err := server.SetupListener(nil)
	if err == nil {
		t.Errorf("Should throw the error")
	}
}

// LISTEN FOR CONNECTION TESTS

func TestListenForConnection(t *testing.T) {
	var server Server = &RealServer{}
	var listener net.Listener = &MockListener{}
	listener.Close()
	listener.Addr()
	_, err := server.ListenForConnection(listener)
	if err == nil {
		t.Errorf("Mock connection, accept should yield an error")
	}
}

// MAKE CONNECTION READER TESTS

func TestMakeConnectionReader(t *testing.T) {
	var server Server = &RealServer{}
	conn, _ := net.Dial("tcp", "1.2.3.4")
	actual := server.MakeConnectionReader(&conn)
	if actual == nil {
		t.Errorf("Expected %#v not to be nil", actual)
	}
}

// COMMAND HANDLER TESTS

func TestCommandGetInvalidID(t *testing.T) {
	sender := &network.MockSender{}
	api := api_p.API{Sender: sender}
	var server Server = &RealServer{api: api}
	expected := "Error: Invalid Kademlia ID;"
	actual := server.CommandHandler([]string{"get", "123"})
	if expected != actual {
		t.Errorf("Expected %v got %v", expected, actual)
	}
}

func TestCommandFindValueFailed(t *testing.T) {
	sender := &network.MockSender{FindValueErr: errors.New("TestFindValueFailed")}
	api := api_p.API{Sender: sender}
	var server Server = &RealServer{api: api}
	actual := server.CommandHandler(
		[]string{"get", "0000000000000000000000000000000000000000"},
	)
	expected := "TestFindValueFailed;"
	if expected != actual {
		t.Errorf("Expected %v got %v", expected, actual)
	}
}

func TestCommandFindValueSuccessful(t *testing.T) {
	sender := &network.MockSender{FindValueResponse: "TestFindValueSuccessful"}
	api := api_p.API{Sender: sender}
	var server Server = &RealServer{api: api}
	actual := server.CommandHandler(
		[]string{"get", "0000000000000000000000000000000000000000"},
	)
	expected := "Value: TestFindValueSuccessful;"
	if expected != actual {
		t.Errorf("Expected %v got %v", expected, actual)
	}
}

func TestCommandStore(t *testing.T) {
	str := "TestFindValueSuccessful"
	key := hashing.NewKademliaID(str)
	sender := &network.MockSender{
		StoreSent: 2,
		FindNodeResponse: [constants.REPLICATION_FACTOR]*nodeutils.Node{
			nil,
			nil,
			nil,
		},
	}
	api := api_p.API{Sender: sender}
	var server Server = &RealServer{api: api}
	actual := server.CommandHandler(
		[]string{"put", str},
	)
	expected := fmt.Sprintf("Stored at: %v on %v nodes.;", key, 2)
	if actual != expected {
		t.Errorf("Expected %v got %v", expected, actual)
	}
}

func TestCommandPing(t *testing.T) {
	sender := &network.MockSender{
		PingResponse: true,
	}
	api := api_p.API{Sender: sender}
	var server Server = &RealServer{api: api}
	actual := server.CommandHandler(
		[]string{"ping", "1.2.3.4"},
	)
	expected := fmt.Sprintf("Online: %v;", true)
	if actual != expected {
		t.Errorf("Expected %v got %v", expected, actual)
	}
}

// SEND MESSAGE TESTS

func TestSendsTheCorrectMessage(t *testing.T) {
	ch := make(chan []byte)
	var writer io.Writer = &MockWriter{ch: ch}
	var server Server = &RealServer{}
	expected := "abc"
	go server.SendMessage(&writer, expected)
	actual := <-ch
	if string(expected) != string(actual) {
		t.Errorf("Expected %v got %v", expected, actual)
	}
}

func TestDoesNotCrashOnError(t *testing.T) {
	ch := make(chan []byte)
	var writer io.Writer = &MockWriter{
		ch:       ch,
		WriteErr: errors.New("TestDoesNotCrashOnError"),
	}
	var server Server = &RealServer{}
	go server.SendMessage(&writer, "abc")
	<-ch
}

// CLOSE CONNECTION TESTS

func TestCloseConnectionShouldNotClose(t *testing.T) {
	var server Server = &RealServer{}
	var closer io.Closer = &MockCloser{}
	cliChannel := make(chan string)
	server.CloseConnection(&closer, cliChannel, "")
}

func TestCloseConnection(t *testing.T) {
	var server Server = &RealServer{}
	var closer io.Closer = &MockCloser{}
	cliChannel := make(chan string)
	server.CloseConnection(&closer, cliChannel, "Closing connection.;")
}

func TestTerminateNode(t *testing.T) {
	var server Server = &RealServer{}
	var closer io.Closer = &MockCloser{}
	cliChannel := make(chan string)
	go server.CloseConnection(&closer, cliChannel, "Terminating node.;")
	<-cliChannel
}
