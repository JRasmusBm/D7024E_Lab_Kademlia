package network

import (
	"encoding/json"
	"errors"
	"io"
	"sync"
	"testing"
	"utils/hashing"
	nodeutils "utils/node"
	storage_p "utils/storage"
)

func TestDial(t *testing.T) {
	var dialer Dialer = &MockDialer{
		index:       0,
		DialResults: []io.ReadWriter{nil},
		DialErrors:  []error{nil},
	}
	var sender Sender = &RealSender{Dialer: &dialer}
	_, err := sender.Dial(&nodeutils.Node{IP: "0.0.0.0"})
	if err != nil {
		t.Error(err)
	}
}

func TestDialFailed(t *testing.T) {
	var dialer Dialer = &MockDialer{
		index:       0,
		DialResults: []io.ReadWriter{nil},
		DialErrors:  []error{errors.New("TestDialFailed")},
	}
	var sender Sender = &RealSender{Dialer: &dialer}
	_, err := sender.Dial(&nodeutils.Node{IP: "0.0.0.0"})
	if err == nil {
		t.Errorf("Should throw an error")
	}
}

func TestRecursiveLookup(t *testing.T) {
	id1, _ := hashing.ToKademliaID("1111111111111111111111111111111111111111")
	id2, _ := hashing.ToKademliaID("2222222222222222222222222222222222222222")
	node1 := nodeutils.Node{ID: id1, IP: "0.0.0.0"}
	node2 := nodeutils.Node{ID: id2, IP: "0.0.0.0"}
	msg, _ := json.Marshal(
		FindNodeRespMsg{Nodes: node1.String() + " " + node2.String()},
	)
	var dialer Dialer = &MockDialer{
		index: 0,
		DialResults: []io.ReadWriter{
			&MockReadWriter{
				Msg: msg,
			},
			&MockReadWriter{
				Msg: msg,
			},
		},
		DialErrors: []error{nil, nil, nil},
	}
	addNode := make(chan nodeutils.AddNodeOp)
	go func() {
		for {
			<-addNode
		}
	}()
	sender := RealSender{Dialer: &dialer, AddNode: addNode}
	candidates := nodeutils.NodeCandidates{}
	queriedNodes := queriedCandidates{nodes: make([]nodeutils.Node, 0)}
	var wg sync.WaitGroup
	wg.Add(1)
	go sender.recursiveLookup(id1, node1, &candidates, &queriedNodes, &wg)
	wg.Wait()
	expected := node1.String()
	actual := candidates.GetNodes(1)[0].String()
	if expected != actual {
		t.Errorf("Expected %v got %v", expected, actual)
	}
}

func TestRecursiveLookupFailedDecode(t *testing.T) {
	id1, _ := hashing.ToKademliaID("1111111111111111111111111111111111111111")
	node1 := nodeutils.Node{ID: id1, IP: "0.0.0.0"}
	msg := []byte("0000")
	var dialer Dialer = &MockDialer{
		index: 0,
		DialResults: []io.ReadWriter{
			&MockReadWriter{
				Msg: msg,
			},
			&MockReadWriter{
				Msg: msg,
			},
		},
		DialErrors: []error{nil, nil, nil},
	}
	addNode := make(chan nodeutils.AddNodeOp)
	go func() {
		for {
			<-addNode
		}
	}()
	sender := RealSender{Dialer: &dialer, AddNode: addNode}
	candidates := nodeutils.NodeCandidates{}
	queriedNodes := queriedCandidates{nodes: make([]nodeutils.Node, 0)}
	var wg sync.WaitGroup
	wg.Add(1)
	go sender.recursiveLookup(id1, node1, &candidates, &queriedNodes, &wg)
	wg.Wait()
	if candidates.Len() != 0 {
		t.Errorf("Expected to find no nodes")
	}
}

func TestPing(t *testing.T) {
	id1, _ := hashing.ToKademliaID("1111111111111111111111111111111111111111")
	node1 := nodeutils.Node{ID: id1, IP: "0.0.0.0"}
	msg, _ := json.Marshal(
		PingMsg{Msg: "PONG"},
	)
	var dialer Dialer = &MockDialer{
		index: 0,
		DialResults: []io.ReadWriter{&MockReadWriter{
			Msg: msg,
		}},
		DialErrors: []error{nil},
	}
	addNode := make(chan nodeutils.AddNodeOp)
	sender := RealSender{Dialer: &dialer, AddNode: addNode}
	ch := make(chan bool)
	errCh := make(chan error)
	go func() {
		for {
			<-addNode
		}
	}()
	go sender.Ping(&node1, ch, errCh)
	var err error
	var actual bool
	select {
	case err = <-errCh:
		// do nothing
	case actual = <-ch:
		// do nothing
	}
	expected := true
	if err != nil {
		t.Errorf(err.Error())
	}
	if expected != actual {
		t.Errorf("Expected %v got %v", expected, actual)
	}
}

func TestPingDialErrors(t *testing.T) {
	id1, _ := hashing.ToKademliaID("1111111111111111111111111111111111111111")
	node1 := nodeutils.Node{ID: id1, IP: "0.0.0.0"}
	msg := []byte("0000")
	var dialer Dialer = &MockDialer{
		index: 0,
		DialResults: []io.ReadWriter{&MockReadWriter{
			Msg:      msg,
			ReadErr:  errors.New("TestPingDialErrors"),
			WriteErr: errors.New("TestPingDialErrors"),
		}},
		DialErrors: []error{nil},
	}
	addNode := make(chan nodeutils.AddNodeOp)
	sender := RealSender{Dialer: &dialer, AddNode: addNode}
	ch := make(chan bool)
	errCh := make(chan error)
	go func() {
		for {
			<-addNode
		}
	}()
	go sender.Ping(&node1, ch, errCh)
	var err error
	select {
	case err = <-errCh:
		// do nothing
	case <-ch:
		// do nothing
	}
	if err == nil {
		t.Errorf("Expected to throw error")
	}
}

func TestPingFailConn(t *testing.T) {
	id1, _ := hashing.ToKademliaID("1111111111111111111111111111111111111111")
	node1 := nodeutils.Node{ID: id1, IP: "0.0.0.0"}
	var dialer Dialer = &MockDialer{
		index:       0,
		DialResults: []io.ReadWriter{nil},
		DialErrors:  []error{errors.New("TestDialFailed")},
	}
	addNode := make(chan nodeutils.AddNodeOp)
	sender := RealSender{Dialer: &dialer, AddNode: addNode}
	ch := make(chan bool)
	errCh := make(chan error)
	go func() {
		for {
			<-addNode
		}
	}()
	go sender.Ping(&node1, ch, errCh)
	var err error
	var actual bool
	select {
	case err = <-errCh:
		// do nothing
	case actual = <-ch:
		// do nothing
	}
	expected := false
	if err == nil {
		t.Errorf(err.Error())
	}
	if expected != actual {
		t.Errorf("Expected %v got %v", expected, actual)
	}
}

func TestStore(t *testing.T) {
	id1, _ := hashing.ToKademliaID("1111111111111111111111111111111111111111")
	id2, _ := hashing.ToKademliaID("2222222222222222222222222222222222222222")
	node1 := nodeutils.Node{ID: id1, IP: "0.0.0.0"}
	node2 := nodeutils.Node{ID: id2, IP: "0.0.0.0"}
	findNodeMsg, _ := json.Marshal(
		FindNodeRespMsg{Nodes: node1.String() + " " + node2.String()},
	)
	var dialer Dialer = &MockDialer{
		index: 0,
		DialResults: []io.ReadWriter{
			&MockReadWriter{
				Msg: findNodeMsg,
			},
			&MockReadWriter{
				Msg: findNodeMsg,
			},
			&MockReadWriter{
				Msg: findNodeMsg,
			},
			&MockReadWriter{
				Msg: findNodeMsg,
			},
		},
		DialErrors: []error{nil, nil, nil, nil},
	}
	addNode := make(chan nodeutils.AddNodeOp)
	findClosestNodes := make(chan nodeutils.FindClosestNodesOp, 1000)
	go func() {
		message := <-findClosestNodes
		message.Resp <- []*nodeutils.Node{&node1}
	}()
	go func() {
		for {
			<-addNode
		}
	}()
	sender := RealSender{Dialer: &dialer,
		AddNode:          addNode,
		FindClosestNodes: findClosestNodes,
	}
	content := "abc"
	sentCh := make(chan int)
	go sender.Store(content, sentCh)
	expected := 2
	actual := <-sentCh
	if expected != actual {
		t.Errorf("Expected %v got %v", expected, actual)
	}
}

func TestStoreDialErrors(t *testing.T) {
	id1, _ := hashing.ToKademliaID("1111111111111111111111111111111111111111")
	id2, _ := hashing.ToKademliaID("2222222222222222222222222222222222222222")
	node1 := nodeutils.Node{ID: id1, IP: "0.0.0.0"}
	node2 := nodeutils.Node{ID: id2, IP: "0.0.0.0"}
	findNodeMsg, _ := json.Marshal(
		FindNodeRespMsg{Nodes: node1.String() + " " + node2.String()},
	)
	var dialer Dialer = &MockDialer{
		index: 0,
		DialResults: []io.ReadWriter{
			&MockReadWriter{
				Msg: findNodeMsg,
			},
			nil,
			nil,
			nil,
		},
		DialErrors: []error{
			nil,
			errors.New("TestStoreDialErrors"),
			errors.New("TestStoreDialErrors"),
			errors.New("TestStoreDialErrors"),
		},
	}
	addNode := make(chan nodeutils.AddNodeOp)
	findClosestNodes := make(chan nodeutils.FindClosestNodesOp, 1000)
	go func() {
		message := <-findClosestNodes
		message.Resp <- []*nodeutils.Node{&node1}
	}()
	go func() {
		for {
			<-addNode
		}
	}()
	sender := RealSender{Dialer: &dialer,
		AddNode:          addNode,
		FindClosestNodes: findClosestNodes,
	}
	content := "abc"
	sentCh := make(chan int)
	go sender.Store(content, sentCh)
	expected := 0
	actual := <-sentCh
	if expected != actual {
		t.Errorf("Expected %v got %v", expected, actual)
	}
}

func TestLookUpValue(t *testing.T) {
	id1, _ := hashing.ToKademliaID("1111111111111111111111111111111111111111")
	id2, _ := hashing.ToKademliaID("2222222222222222222222222222222222222222")
	node1 := nodeutils.Node{ID: id1, IP: "0.0.0.0"}
	node2 := nodeutils.Node{ID: id2, IP: "0.0.0.0"}
	expected := "abc"
	lookupMsg, _ := json.Marshal(
		FindNodeRespMsg{
			Nodes: nodeutils.ToStrings([]*nodeutils.Node{&node1, &node2}),
		},
	)
	lookupValueMsg, _ := json.Marshal(
		FindValueRespMsg{Content: expected, Nodes: ""},
	)
	var storage storage_p.Storage = &MockStorage{ReadErr: errors.New("LookUpValue")}
	var dialer Dialer = &MockDialer{
		index: 0,
		DialResults: []io.ReadWriter{
			&MockReadWriter{
				Msg: lookupMsg,
			},
			&MockReadWriter{
				Msg: lookupValueMsg,
			},
			&MockReadWriter{
				Msg: lookupValueMsg,
			},
			&MockReadWriter{
				Msg: lookupValueMsg,
			},
			&MockReadWriter{
				Msg: lookupValueMsg,
			},
			&MockReadWriter{
				Msg: lookupValueMsg,
			},
		},
		DialErrors: []error{nil, nil, nil, nil, nil},
	}
	addNode := make(chan nodeutils.AddNodeOp)
	findClosestNodes := make(chan nodeutils.FindClosestNodesOp, 1000)
	go func() {
		for {
			message := <-findClosestNodes
			message.Resp <- []*nodeutils.Node{&node1, &node2}
		}
	}()
	go func() {
		for {
			<-addNode
		}
	}()
	sender := RealSender{
		Dialer:           &dialer,
		Storage:          &storage,
		AddNode:          addNode,
		FindClosestNodes: findClosestNodes,
	}
	key := hashing.NewKademliaID(expected)
	actual := sender.LookUpValue(key)
	if expected != actual {
		t.Errorf("Expected %v got %v", expected, actual)
	}
}

func TestLookUpValueNodes(t *testing.T) {
	id1, _ := hashing.ToKademliaID("1111111111111111111111111111111111111111")
	id2, _ := hashing.ToKademliaID("2222222222222222222222222222222222222222")
	node1 := nodeutils.Node{ID: id1, IP: "0.0.0.0"}
	node2 := nodeutils.Node{ID: id2, IP: "0.0.0.0"}
	expected := "abc"
	lookupMsg, _ := json.Marshal(
		FindNodeRespMsg{
			Nodes: nodeutils.ToStrings([]*nodeutils.Node{&node1, &node2}),
		},
	)
	lookupValueMsg, _ := json.Marshal(
		FindValueRespMsg{Content: "", Nodes: nodeutils.ToStrings([]*nodeutils.Node{&node1, &node2})},
	)
	var storage storage_p.Storage = &MockStorage{ReadErr: errors.New("LookUpValue")}
	var dialer Dialer = &MockDialer{
		index: 0,
		DialResults: []io.ReadWriter{
			&MockReadWriter{
				Msg: lookupMsg,
			},
			&MockReadWriter{
				Msg: lookupValueMsg,
			},
			&MockReadWriter{
				Msg: lookupValueMsg,
			},
			&MockReadWriter{
				Msg: lookupValueMsg,
			},
			&MockReadWriter{
				Msg: lookupValueMsg,
			},
			&MockReadWriter{
				Msg: lookupValueMsg,
			},
		},
		DialErrors: []error{nil, nil, nil, nil, nil},
	}
	addNode := make(chan nodeutils.AddNodeOp)
	findClosestNodes := make(chan nodeutils.FindClosestNodesOp, 1000)
	go func() {
		for {
			message := <-findClosestNodes
			message.Resp <- []*nodeutils.Node{&node1, &node2}
		}
	}()
	go func() {
		for {
			<-addNode
		}
	}()
	sender := RealSender{
		Dialer:           &dialer,
		Storage:          &storage,
		AddNode:          addNode,
		FindClosestNodes: findClosestNodes,
	}
	key := hashing.NewKademliaID(expected)
	actual := sender.LookUpValue(key)
	if "" != actual {
		t.Errorf("Should not return value when not found, actual: %v", actual)
	}
}

func TestLookUpValueDialError(t *testing.T) {
	id1, _ := hashing.ToKademliaID("1111111111111111111111111111111111111111")
	id2, _ := hashing.ToKademliaID("2222222222222222222222222222222222222222")
	node1 := nodeutils.Node{ID: id1, IP: "0.0.0.0"}
	node2 := nodeutils.Node{ID: id2, IP: "0.0.0.0"}
	expected := "abc"
	lookupMsg, _ := json.Marshal(
		FindNodeRespMsg{
			Nodes: nodeutils.ToStrings([]*nodeutils.Node{&node1, &node2}),
		},
	)
	lookupValueMsg, _ := json.Marshal(
		FindValueRespMsg{Content: expected, Nodes: ""},
	)
	var storage storage_p.Storage = &MockStorage{ReadErr: errors.New("LookUpValue")}
	var dialer Dialer = &MockDialer{
		index: 0,
		DialResults: []io.ReadWriter{
			&MockReadWriter{
				Msg: lookupMsg,
			},
			&MockReadWriter{
				Msg: lookupValueMsg,
			},
			&MockReadWriter{
				Msg: lookupValueMsg,
			},
			&MockReadWriter{
				Msg: lookupValueMsg,
			},
			&MockReadWriter{
				Msg: lookupValueMsg,
			},
			&MockReadWriter{
				Msg: lookupValueMsg,
			},
		},
		DialErrors: []error{
			nil,
			errors.New("TestLookUpValueDialError"),
			errors.New("TestLookUpValueDialError"),
			errors.New("TestLookUpValueDialError"),
			errors.New("TestLookUpValueDialError"),
		},
	}
	addNode := make(chan nodeutils.AddNodeOp)
	findClosestNodes := make(chan nodeutils.FindClosestNodesOp, 1000)
	go func() {
		for {
			message := <-findClosestNodes
			message.Resp <- []*nodeutils.Node{&node1, &node2}
		}
	}()
	go func() {
		for {
			<-addNode
		}
	}()
	sender := RealSender{
		Dialer:           &dialer,
		Storage:          &storage,
		AddNode:          addNode,
		FindClosestNodes: findClosestNodes,
	}
	key := hashing.NewKademliaID(expected)
	actual := sender.LookUpValue(key)
	if "" != actual {
		t.Errorf("Expected empty string got %v", actual)
	}
}

func TestLookUpValueDecodeError(t *testing.T) {
	id1, _ := hashing.ToKademliaID("1111111111111111111111111111111111111111")
	id2, _ := hashing.ToKademliaID("2222222222222222222222222222222222222222")
	node1 := nodeutils.Node{ID: id1, IP: "0.0.0.0"}
	node2 := nodeutils.Node{ID: id2, IP: "0.0.0.0"}
	expected := "abc"
	lookupMsg, _ := json.Marshal(
		FindNodeRespMsg{
			Nodes: nodeutils.ToStrings([]*nodeutils.Node{&node1, &node2}),
		},
	)
	lookupValueMsg := []byte("0000")
	var storage storage_p.Storage = &MockStorage{ReadErr: errors.New("LookUpValue")}
	var dialer Dialer = &MockDialer{
		index: 0,
		DialResults: []io.ReadWriter{
			&MockReadWriter{
				Msg: lookupMsg,
			},
			&MockReadWriter{
				Msg: lookupValueMsg,
			},
			&MockReadWriter{
				Msg: lookupValueMsg,
			},
			&MockReadWriter{
				Msg: lookupValueMsg,
			},
			&MockReadWriter{
				Msg: lookupValueMsg,
			},
			&MockReadWriter{
				Msg: lookupValueMsg,
			},
		},
		DialErrors: []error{
			nil,
			nil,
			nil,
			nil,
			nil,
		},
	}
	addNode := make(chan nodeutils.AddNodeOp)
	findClosestNodes := make(chan nodeutils.FindClosestNodesOp, 1000)
	go func() {
		for {
			message := <-findClosestNodes
			message.Resp <- []*nodeutils.Node{&node1, &node2}
		}
	}()
	go func() {
		for {
			<-addNode
		}
	}()
	sender := RealSender{
		Dialer:           &dialer,
		Storage:          &storage,
		AddNode:          addNode,
		FindClosestNodes: findClosestNodes,
	}
	key := hashing.NewKademliaID(expected)
	actual := sender.LookUpValue(key)
	if "" != actual {
		t.Errorf("Expected empty string got %v", actual)
	}
}

func TestLookUpValueEmptyString(t *testing.T) {
	id1, _ := hashing.ToKademliaID("1111111111111111111111111111111111111111")
	id2, _ := hashing.ToKademliaID("2222222222222222222222222222222222222222")
	node1 := nodeutils.Node{ID: id1, IP: "0.0.0.0"}
	node2 := nodeutils.Node{ID: id2, IP: "0.0.0.0"}
	expected := ""
	lookupMsg, _ := json.Marshal(
		FindNodeRespMsg{
			Nodes: nodeutils.ToStrings([]*nodeutils.Node{&node1, &node2}),
		},
	)
	lookupValueMsg, _ := json.Marshal(
		FindValueRespMsg{Content: expected, Nodes: ""},
	)
	var storage storage_p.Storage = &MockStorage{ReadErr: errors.New("LookUpValue")}
	var dialer Dialer = &MockDialer{
		index: 0,
		DialResults: []io.ReadWriter{
			&MockReadWriter{
				Msg: lookupMsg,
			},
			&MockReadWriter{
				Msg: lookupValueMsg,
			},
			&MockReadWriter{
				Msg: lookupValueMsg,
			},
			&MockReadWriter{
				Msg: lookupValueMsg,
			},
			&MockReadWriter{
				Msg: lookupValueMsg,
			},
			&MockReadWriter{
				Msg: lookupValueMsg,
			},
		},
		DialErrors: []error{nil, nil, nil, nil, nil},
	}
	addNode := make(chan nodeutils.AddNodeOp)
	findClosestNodes := make(chan nodeutils.FindClosestNodesOp, 1000)
	go func() {
		for {
			message := <-findClosestNodes
			message.Resp <- []*nodeutils.Node{&node1, &node2}
		}
	}()
	go func() {
		for {
			<-addNode
		}
	}()
	sender := RealSender{
		Dialer:           &dialer,
		Storage:          &storage,
		AddNode:          addNode,
		FindClosestNodes: findClosestNodes,
	}
	key := hashing.NewKademliaID(expected)
	actual := sender.LookUpValue(key)
	if expected != actual {
		t.Errorf("Expected %v got %v", expected, actual)
	}
}

func TestLookUpValueCaseContent(t *testing.T) {
	id1, _ := hashing.ToKademliaID("1111111111111111111111111111111111111111")
	id2, _ := hashing.ToKademliaID("2222222222222222222222222222222222222222")
	node1 := nodeutils.Node{ID: id1, IP: "0.0.0.0"}
	node2 := nodeutils.Node{ID: id2, IP: "0.0.0.0"}
	expected := "abc"
	lookupMsg, _ := json.Marshal(
		FindNodeRespMsg{
			Nodes: nodeutils.ToStrings([]*nodeutils.Node{&node1, &node2}),
		},
	)
	lookupValueMsg, _ := json.Marshal(
		FindValueRespMsg{Content: expected, Nodes: ""},
	)
	var storage storage_p.Storage = &MockStorage{ReadResult: "abc"}
	var dialer Dialer = &MockDialer{
		index: 0,
		DialResults: []io.ReadWriter{
			&MockReadWriter{
				Msg: lookupMsg,
			},
			&MockReadWriter{
				Msg: lookupValueMsg,
			},
			&MockReadWriter{
				Msg: lookupValueMsg,
			},
			&MockReadWriter{
				Msg: lookupValueMsg,
			},
			&MockReadWriter{
				Msg: lookupValueMsg,
			},
			&MockReadWriter{
				Msg: lookupValueMsg,
			},
		},
		DialErrors: []error{nil, nil, nil, nil, nil},
	}
	addNode := make(chan nodeutils.AddNodeOp)
	findClosestNodes := make(chan nodeutils.FindClosestNodesOp, 1000)
	go func() {
		for {
			message := <-findClosestNodes
			message.Resp <- []*nodeutils.Node{&node1, &node2}
		}
	}()
	go func() {
		for {
			<-addNode
		}
	}()
	sender := RealSender{
		Dialer:           &dialer,
		Storage:          &storage,
		AddNode:          addNode,
		FindClosestNodes: findClosestNodes,
	}
	key := hashing.NewKademliaID(expected)
	actual := sender.LookUpValue(key)
	if expected != actual {
		t.Errorf("Expected %v got %v", expected, actual)
	}
}

func TestJoin(t *testing.T) {
	id1, _ := hashing.ToKademliaID("1111111111111111111111111111111111111111")
	node1 := nodeutils.Node{ID: id1, IP: "0.0.0.0"}
	key, _ := hashing.ToKademliaID("0000000000000000000000000000000000000000")
	joinRespMsg, _ := json.Marshal(
		JoinRespMsg{Success: true, ID: key.String()},
	)
	lookupMsg, _ := json.Marshal(
		FindNodeRespMsg{
			Nodes: nodeutils.ToStrings([]*nodeutils.Node{&node1, &node1}),
		},
	)
	var dialer Dialer = &MockDialer{
		index: 0,
		DialResults: []io.ReadWriter{
			&MockReadWriter{
				Msg: joinRespMsg,
			},
			&MockReadWriter{
				Msg: lookupMsg,
			},
		},
		DialErrors: []error{nil, nil, nil, nil, nil},
	}
	addNode := make(chan nodeutils.AddNodeOp)
	findClosestNodes := make(chan nodeutils.FindClosestNodesOp, 1000)
	go func() {
		message := <-findClosestNodes
		message.Resp <- []*nodeutils.Node{&node1}
	}()
	go func() {
		for {
			<-addNode
		}
	}()
	sender := RealSender{
		Me:               &node1,
		FindClosestNodes: findClosestNodes,
		Dialer:           &dialer,
		AddNode:          addNode,
	}
	successCh := make(chan bool)
	errCh := make(chan error)
	go sender.Join("0.0.0.0", successCh, errCh)
	var actual bool
	var err error
	expected := true
	select {
	case actual = <-successCh:
		// Do nothing
	case err = <-errCh:
		// Do nothing
	}
	if err != nil {
		t.Errorf(err.Error())
	}
	if expected != actual {
		t.Errorf("Expected %v got %v", expected, actual)
	}
}

func TestJoinDialError(t *testing.T) {
	id1, _ := hashing.ToKademliaID("1111111111111111111111111111111111111111")
	node1 := nodeutils.Node{ID: id1, IP: "0.0.0.0"}
	key, _ := hashing.ToKademliaID("0000000000000000000000000000000000000000")
	joinRespMsg, _ := json.Marshal(
		JoinRespMsg{Success: true, ID: key.String()},
	)
	lookupMsg, _ := json.Marshal(
		FindNodeRespMsg{
			Nodes: nodeutils.ToStrings([]*nodeutils.Node{&node1, &node1}),
		},
	)
	var dialer Dialer = &MockDialer{
		index: 0,
		DialResults: []io.ReadWriter{
			&MockReadWriter{
				Msg: joinRespMsg,
			},
			&MockReadWriter{
				Msg: lookupMsg,
			},
		},
		DialErrors: []error{
			errors.New("TestJoinDialError"),
			errors.New("TestJoinDialError"),
			errors.New("TestJoinDialError"),
			errors.New("TestJoinDialError"),
			errors.New("TestJoinDialError"),
		},
	}
	addNode := make(chan nodeutils.AddNodeOp)
	findClosestNodes := make(chan nodeutils.FindClosestNodesOp, 1000)
	go func() {
		message := <-findClosestNodes
		message.Resp <- []*nodeutils.Node{&node1}
	}()
	go func() {
		for {
			<-addNode
		}
	}()
	sender := RealSender{
		Me:               &node1,
		FindClosestNodes: findClosestNodes,
		Dialer:           &dialer,
		AddNode:          addNode,
	}
	successCh := make(chan bool)
	errCh := make(chan error)
	go sender.Join("0.0.0.0", successCh, errCh)
	var actual bool
	var err error
	expected := false
	select {
	case actual = <-successCh:
		// Do nothing
	case err = <-errCh:
		// Do nothing
	}
	if err == nil {
		t.Errorf("Should return error")
	}
	if expected != actual {
		t.Errorf("Expected %v got %v", expected, actual)
	}
}

func TestJoinDecodeError(t *testing.T) {
	id1, _ := hashing.ToKademliaID("1111111111111111111111111111111111111111")
	node1 := nodeutils.Node{ID: id1, IP: "0.0.0.0"}
	joinRespMsg := []byte("0000")
	lookupMsg, _ := json.Marshal(
		FindNodeRespMsg{
			Nodes: nodeutils.ToStrings([]*nodeutils.Node{&node1, &node1}),
		},
	)
	var dialer Dialer = &MockDialer{
		index: 0,
		DialResults: []io.ReadWriter{
			&MockReadWriter{
				Msg: joinRespMsg,
			},
			&MockReadWriter{
				Msg: lookupMsg,
			},
		},
		DialErrors: []error{
			nil,
			nil,
			nil,
			nil,
			nil,
		},
	}
	addNode := make(chan nodeutils.AddNodeOp)
	findClosestNodes := make(chan nodeutils.FindClosestNodesOp, 1000)
	go func() {
		message := <-findClosestNodes
		message.Resp <- []*nodeutils.Node{&node1}
	}()
	go func() {
		for {
			<-addNode
		}
	}()
	sender := RealSender{
		Me:               &node1,
		FindClosestNodes: findClosestNodes,
		Dialer:           &dialer,
		AddNode:          addNode,
	}
	successCh := make(chan bool)
	errCh := make(chan error)
	go sender.Join("0.0.0.0", successCh, errCh)
	var actual bool
	var err error
	expected := false
	select {
	case actual = <-successCh:
		// Do nothing
	case err = <-errCh:
		// Do nothing
	}
	if err != nil {
		t.Errorf("Should return error")
	}
	if expected != actual {
		t.Errorf("Expected %v got %v", expected, actual)
	}
}

func TestMockStorageWrite(t *testing.T) {
	storage := MockStorage{}
	key := hashing.NewRandomKademliaID()
	storage.Write(key.String(), "abc")
}
