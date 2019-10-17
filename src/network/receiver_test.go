package network

import (
	"encoding/json"
	"fmt"
	"io"
	"testing"
	"time"
	"utils/hashing"
	nodeutils "utils/node"
	storage_p "utils/storage"
)

func TestServerDecodeError(t *testing.T) {
	conn := MockReadWriter{}
	listener := MockListener{AcceptResult: conn}
	var sender Sender = &MockSender{}
	var storage storage_p.Storage = &MockStorage{ReadResult: ""}
	receiver := RealReceiver{
		Sender:   sender,
		Listener: listener,
		Storage:  &storage,
	}
	go receiver.Server()
	time.Sleep(300 * time.Millisecond)
}

func TestServerPing(t *testing.T) {
	pingMsg, _ := json.Marshal(
		Message{RPC: "PING", Msg: PingMsg{Msg: "PING"}},
	)
	conn := MockReadWriter{
		Msg: pingMsg,
	}
	listener := MockListener{AcceptResult: conn}
	var sender Sender = &MockSender{}
	var storage storage_p.Storage = &MockStorage{ReadResult: ""}
	receiver := RealReceiver{
		Sender:   sender,
		Listener: listener,
		Storage:  &storage,
	}
	go receiver.Server()
	time.Sleep(300 * time.Millisecond)
}

func TestServerFindNode(t *testing.T) {
	findNodeMsg, _ := json.Marshal(
		Message{
			RPC: "FIND_NODE",
			Msg: FindNodeMsg{ID: "0000000000000000000000000000000000000000"},
		},
	)
	fmt.Printf(string(findNodeMsg))
	conn := MockReadWriter{
		Msg: findNodeMsg,
	}
	listener := MockListener{AcceptResult: conn}
	var sender Sender = &MockSender{}
	var storage storage_p.Storage = &MockStorage{ReadResult: ""}
	receiver := RealReceiver{
		Sender:   sender,
		Listener: listener,
		Storage:  &storage,
	}
	go receiver.Server()
	time.Sleep(300 * time.Millisecond)
}

func TestServerStore(t *testing.T) {
	findNodeMsg, _ := json.Marshal(
		Message{
			RPC: "STORE",
			Msg: StoreMsg{Data: "Abc"},
		},
	)
	fmt.Printf(string(findNodeMsg))
	conn := MockReadWriter{
		Msg: findNodeMsg,
	}
	listener := MockListener{AcceptResult: conn}
	var sender Sender = &MockSender{}
	var storage storage_p.Storage = &MockStorage{ReadResult: ""}
	receiver := RealReceiver{
		Sender:   sender,
		Listener: listener,
		Storage:  &storage,
	}
	go receiver.Server()
	time.Sleep(300 * time.Millisecond)
}

func TestServerJoin(t *testing.T) {
	findNodeMsg, _ := json.Marshal(
		Message{
			RPC: "JOIN",
			Msg: JoinMsg{Msg: "Hello"},
		},
	)
	fmt.Printf(string(findNodeMsg))
	conn := MockReadWriter{
		Msg: findNodeMsg,
	}
	listener := MockListener{AcceptResult: conn}
	var sender Sender = &MockSender{}
	var storage storage_p.Storage = &MockStorage{ReadResult: ""}
	receiver := RealReceiver{
		Sender:   sender,
		Listener: listener,
		Storage:  &storage,
	}
	go receiver.Server()
	time.Sleep(300 * time.Millisecond)
}

func TestServerFindValue(t *testing.T) {
	findNodeMsg, _ := json.Marshal(
		Message{
			RPC: "FIND_VALUE",
			Msg: FindValueMsg{Key: "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"},
		},
	)
	fmt.Printf(string(findNodeMsg))
	conn := MockReadWriter{
		Msg: findNodeMsg,
	}
	listener := MockListener{AcceptResult: conn}
	var sender Sender = &MockSender{}
	var storage storage_p.Storage = &MockStorage{ReadResult: ""}
	receiver := RealReceiver{
		Sender:   sender,
		Listener: listener,
		Storage:  &storage,
	}
	go receiver.Server()
	time.Sleep(300 * time.Millisecond)
}

func TestJoinReply(t *testing.T) {
	id1, _ := hashing.ToKademliaID("1111111111111111111111111111111111111111")
	node1 := nodeutils.Node{ID: id1, IP: "0.0.0.0"}
	msg := Message{
		RPC: "JOIN",
		Msg: JoinMsg{Msg: "Hello"},
	}
	var storage storage_p.Storage = &MockStorage{ReadResult: ""}
	addNode := make(chan nodeutils.AddNodeOp)
	go func() {
		for {
			<-addNode
		}
	}()
	receiver := RealReceiver{
		Storage: &storage,
		AddNode: addNode,
		Me:      &node1,
	}
	ch := make(chan []byte)
	var conn io.ReadWriter = &MockReadWriter{
		WriteCh: ch,
		i:       0,
	}
	go receiver.JoinReply(msg, conn)
	actual := <-ch
	data, _ := json.Marshal(
    JoinRespMsg{Success: true, ID: id1.String(), IP: "0.0.0.0"},
	)
	expected := string(data) + "\n"
	if expected != string(actual) {
		t.Errorf("Expected '%v' got '%v'", string(expected), string(actual))
	}
}

func TestFindNodeReply(t *testing.T) {
	id1, _ := hashing.ToKademliaID("1111111111111111111111111111111111111111")
	node1 := nodeutils.Node{ID: id1, IP: "0.0.0.0"}
	msg := Message{
		RPC: "FIND_NODE",
		Msg: FindNodeMsg{ID: id1.String()},
	}
	var storage storage_p.Storage = &MockStorage{ReadResult: ""}
	addNode := make(chan nodeutils.AddNodeOp)
	go func() {
		for {
			<-addNode
		}
	}()

  findClosestNodes := make(chan nodeutils.FindClosestNodesOp, 1000)
  go func() {
    message := <-findClosestNodes
    message.Resp <- []nodeutils.Node{node1}
  }()
	receiver := RealReceiver{
		Storage: &storage,
    AddNode: addNode,
		FindClosestNodes: findClosestNodes,
		Me:      &node1,
	}

	var conn io.ReadWriter = &MockReadWriter{}
	go receiver.FindNodeReply(msg, conn)

}
