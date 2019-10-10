package cli

import (
	"errors"
	"testing"
	"time"
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
