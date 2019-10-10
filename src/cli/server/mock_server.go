package cli

import (
	"net"
)

type MockServer struct {
	SetupListenerResult        net.Listener
	SetupListenerErr           error
	ListenForConnectionResult  net.Conn
	ListenForConnectionErr     error
	MakeConnectionReaderResult *Reader
	ListenToClientResult       string
	MessageParserResult        []string
	CommandHandlerResult       string
}

func (m *MockServer) SetupListener(network *Network) (net.Listener, error) {
	if m.SetupListenerErr != nil {
		return nil, m.SetupListenerErr
	}
	return m.SetupListenerResult, nil
}

func (m *MockServer) ListenForConnection(listener net.Listener) (net.Conn, error) {
	if m.ListenForConnectionErr != nil {
		return nil, m.ListenForConnectionErr
	}
	return m.ListenForConnectionResult, nil
}

func (m *MockServer) MakeConnectionReader(conn *net.Conn) *Reader {
	return m.MakeConnectionReaderResult
}

func (m *MockServer) ListenToClient(reader *Reader) string {
	return m.ListenToClientResult
}

func (m *MockServer) MessageParser(incomingMessage string) []string {
	return m.MessageParserResult
}

func (m *MockServer) CommandHandler(parsedMessage []string) string {
	return m.CommandHandlerResult
}

func (m *MockServer) SendMessage(conn *net.Conn, cliChannel chan string, responseMessage string) {
	return
}
