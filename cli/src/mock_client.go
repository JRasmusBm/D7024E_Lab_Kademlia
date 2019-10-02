package main

import (
	"io"
	"net"
)

type MockClient struct {
	SetupConnectionResult      net.Conn
	SetupConnectionErr         error
	ConnectionValidResult      bool
	MakeConnectionReaderResult *Reader
	GetMessageFromUserResult   string
	GetMessageFromUserErr      error
	HandleMessageResult        string
	HandleMessageErr           error
}

func (m *MockClient) SetupConnection(reader *Reader, network *Network) (net.Conn, error) {
	if m.SetupConnectionErr != nil {
		return nil, m.SetupConnectionErr
	}
	return m.SetupConnectionResult, nil
}

func (m *MockClient) ConnectionValid(conn *io.Writer, reader *Reader) bool {
	return m.ConnectionValidResult
}

func (m *MockClient) MakeConnectionReader(conn *net.Conn) *Reader {
	return m.MakeConnectionReaderResult
}

func (m *MockClient) ListenToServer(reader *Reader) {
	return
}

func (m *MockClient) GetMessageFromUser(reader *Reader) (string, error) {
	if m.GetMessageFromUserErr != nil {
		return "", m.GetMessageFromUserErr
	}
	return m.GetMessageFromUserResult, nil
}

func (m *MockClient) HandleMessage(message string, fileReader *FileReader) (string, error) {
	if m.HandleMessageErr != nil {
		return "", m.HandleMessageErr
	}
	return m.HandleMessageResult, nil
}

func (m *MockClient) SendMessage(conn *io.Writer, rpc string) {
	return
}
