package network

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"utils/constants"
	hashing "utils/hashing"
	nodeutils "utils/node"
	"utils/storage"
)

type Dialer interface {
	DialIP(ip string) (io.ReadWriter, error)
}

type Sender interface {
	Dial(node nodeutils.Node) (io.ReadWriter, error)
	Ping(node nodeutils.Node, ch chan bool, errCh chan error)
	Store(content string, ch chan int)
	FindNode(id *hashing.KademliaID, node nodeutils.Node, ch chan []nodeutils.Node, errCh chan error)
	FindValue(node nodeutils.Node, key *hashing.KademliaID, successCh chan string, closerCh chan [constants.CLOSESTNODES]nodeutils.Node, errCh chan error)
	Join(ip string, ch chan bool, errCh chan error)
	LookUp(id *hashing.KademliaID) []nodeutils.Node
	LookUpValue(key *hashing.KademliaID) string
}

type RealSender struct {
	AddNode          chan nodeutils.AddNodeOp
	FindClosestNodes chan nodeutils.FindClosestNodesOp
	Storage          *storage.Storage
	Me               *nodeutils.Node
	Dialer           *Dialer
}

type queriedCandidates struct {
	sync.RWMutex
	nodes []nodeutils.Node
}

func (sender RealSender) Dial(node nodeutils.Node) (io.ReadWriter, error) {
	return (*sender.Dialer).DialIP(node.IP)
}

func (sender RealSender) LookUp(id *hashing.KademliaID) []nodeutils.Node {
	resp := make(chan []nodeutils.Node)
	sender.FindClosestNodes <- nodeutils.FindClosestNodesOp{Target: id, Count: constants.CONCURRENCY_PARAM, Resp: resp}

	askNodes := <-resp
	fmt.Printf("\nAskNodes: %v\n", askNodes)

	candidates := nodeutils.NodeCandidates{}

	queriedNodes := queriedCandidates{nodes: make([]nodeutils.Node, 0)}

	var wg sync.WaitGroup

	for _, node := range askNodes {
		wg.Add(1)
		go sender.recursiveLookup(id, node, &candidates, &queriedNodes, &wg)
	}

	// Wait for each goroutine to finish
	wg.Wait()
	// fmt.Printf("%v\n", candidates.GetNodes(constants.CLOSESTNODES))

	fmt.Printf("\n LookUp, Candidate Nodes: %v\n", candidates.Nodes)
	return candidates.GetNodes(constants.CLOSESTNODES)
}

func (sender RealSender) recursiveLookup(id *hashing.KademliaID, node nodeutils.Node, candidates *nodeutils.NodeCandidates, queriedNodes *queriedCandidates, wg *sync.WaitGroup) {
	defer wg.Done() // Informs main goroutine that we're done after execution

	ch := make(chan []nodeutils.Node)
	errCh := make(chan error)
	go sender.FindNode(id, node, ch, errCh)
	select {
	case nodes := <-ch:
		// Calculate distance from id for each node, add to candidates and sort.
		var uniqueNodes []nodeutils.Node
		for _, foundNode := range nodes {
			// If we already have the node in our list of candidates, ignore it.
			// Have to use the mutex here in order to properly avoid duplicates.
			candidates.RLock()
			if nodeutils.NodeInArr(foundNode, candidates.Nodes) {
				candidates.RUnlock()
				continue
			}
			candidates.RUnlock()

			foundNode.CalcDistance(id)
			uniqueNodes = append(uniqueNodes, foundNode)
		}
		candidates.Append(uniqueNodes)
		candidates.Sort()

	case err := <- errCh:
		// Connection to node failed, do nothing.
		fmt.Printf("ALERT! Error in FIND_NODE: %s", err.Error())
	}

	queriedNodes.Lock()
	queriedNodes.nodes = append(queriedNodes.nodes, node)
	qNodes := queriedNodes.nodes
	queriedNodes.Unlock()

	// Go through current k-closest nodes and query the first one that hasn't been queried yet.
	kNodes := candidates.GetNodes(constants.CLOSESTNODES)
	fmt.Printf("\n Recursive Lookup, Candidate Nodes: %v\n", candidates.Nodes)
	for _, kNode := range kNodes {
		if !nodeutils.NodeInArr(kNode, qNodes) {
			wg.Add(1)
			sender.recursiveLookup(id, kNode, candidates, queriedNodes, wg)
			return
		}
	}
}

func (sender RealSender) Ping(node nodeutils.Node, ch chan bool, errCh chan error) {
	conn, err := sender.Dial(node)
	if err != nil {
		errCh <- err
		return
	}

	decoder := json.NewDecoder(conn)

	// Send PING message
	encoder := json.NewEncoder(conn)
	encoder.Encode(Message{RPC: "PING", Msg: PingMsg{Msg: "PING"}})

	// Wait for PONG message
	var msg PingMsg
	err = decoder.Decode(&msg)
	fmt.Printf("\nRESPONSE RECEIVED:\n\tRPC: PING\n\tMsg: %v\n",
		msg.Msg,
	)

	if err != nil {
		errCh <- err
	}

	if &node != nil {
		sender.AddNode <- nodeutils.AddNodeOp{AddedNode: node}
	}
	ch <- (msg.Msg == "PONG")
}

func (sender RealSender) Store(content string, ch chan int) {
	fmt.Printf("\nSENDING STORE!\n")
	var conn io.ReadWriter
	var err error
	sent := 0
	id := hashing.NewKademliaID(content)
	fmt.Printf("\nID: %v\n", id.String())
	nodes := sender.LookUp(id)
	fmt.Printf("\n(STORE)Nodes: %v\n", nodes)
	for _, node := range nodes {
		if &node == nil {
			continue
		}

		conn, err = sender.Dial(node)
		if err != nil {
			continue
		}

		encoder := json.NewEncoder(conn)
		encoder.Encode(Message{RPC: "STORE", Msg: StoreMsg{Data: content}})
		sent += 1
	}
	ch <- sent
}

func (sender RealSender) FindNode(id *hashing.KademliaID, node nodeutils.Node, ch chan []nodeutils.Node, errCh chan error) {
	// fmt.Printf("Finding Node %s", node.ID.String())

	conn, err := sender.Dial(node)
	if err != nil {
		errCh <- err
		return
	}

	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)
	encoder.Encode(Message{RPC: "FIND_NODE", Msg: FindNodeMsg{ID: id.String()}})

	var msg FindNodeRespMsg
	err = decoder.Decode(&msg)
	fmt.Printf("\nRESPONSE RECEIVED:\n\tRPC: FIND_NODE\n\tNodes: %v",
		msg.Nodes,
	)
	if err != nil {
		fmt.Printf(err.Error())
		errCh <- err
		return
	}

	nodes, err := nodeutils.FromStrings(msg.Nodes)
	if err != nil {
		fmt.Printf(err.Error())
		errCh <- err
		return
	}

	// Add all given nodes to routing table
	for _, node := range nodes {
		if &node != nil {
			sender.AddNode <- nodeutils.AddNodeOp{AddedNode: node}
		}
	}

	ch <- nodes
}

func (sender RealSender) FindValue(node nodeutils.Node, key *hashing.KademliaID, successCh chan string, closerCh chan [constants.CLOSESTNODES]nodeutils.Node, errCh chan error) {
	readCh := make(chan string)
	readErrCh := make(chan error)
	go (*sender.Storage).Read(key.String(), readCh, readErrCh)
	select {
	case content := <-readCh:
		successCh <- content
		return
	case <-readErrCh:
		// Do nothing
	}

	conn, err := sender.Dial(node)
	if err != nil {
		errCh <- err
		return
	}

	encoder := json.NewEncoder(conn)
	decoder := json.NewDecoder(conn)

	// Send RPC message
	encoder.Encode(Message{RPC: "FIND_VALUE", Msg: FindValueMsg{Key: key.String()}})

	var msg FindValueRespMsg
	err = decoder.Decode(&msg)
	fmt.Printf("\nRESPONSE RECEIVED:\n\tRPC: FIND_VALUE\n\tContent: %v\n\tNodes: %v",
		msg.Content,
		msg.Nodes,
	)
	if err != nil {
		errCh <- err
		return
	}
	fmt.Printf(msg.Content)

	if msg.Content == "" {
		nodes, _ := nodeutils.FromStrings(msg.Nodes)
		var result [constants.CLOSESTNODES]nodeutils.Node
		for i, node := range nodes {
			result[i] = node
		}
		closerCh <- result
	} else {
		successCh <- msg.Content
	}
}

func (sender RealSender) LookUpValue(key *hashing.KademliaID) string {
	// Get the k-closest nodes.
	kNodes := sender.LookUp(key)
	// fmt.Printf("%#v\n", kNodes)
	var wg sync.WaitGroup
	var content string
	var mutex sync.RWMutex

	// NOTE: Currently this will spawn k goroutines, but there should either be 1 or alpha according to the paper.
	// Therefore we should probably make sure it only spawns alpha goroutines in a similar fashion
	// to how it's done in recursiveLookup.
	for _, node := range kNodes {
		if &node != nil {
			wg.Add(1)
			go func(n nodeutils.Node) {
				defer wg.Done()
				successCh := make(chan string)
				closerCh := make(chan [constants.CLOSESTNODES]nodeutils.Node)
				errCh := make(chan error)

				go sender.FindValue(n, key, successCh, closerCh, errCh)
				select {
				case foundContent := <-successCh:
					mutex.Lock()
					content = foundContent
					mutex.Unlock()
				case <-closerCh:
					// Do nothing, node didn't have content (and we don't recursively find new nodes as this is already done previously).
				case <-errCh:
					// Do nothing, couldn't connect to node (?!)
				}
			}(node)
		}
	}

	wg.Wait()
	return content
}

func (sender RealSender) Join(ip string, ch chan bool, errCh chan error) {
	conn, err := (*sender.Dialer).DialIP(ip)
	if err != nil {
		errCh <- err
	}

	decoder := json.NewDecoder(conn)
	// encoder := json.NewEncoder(conn)
	joinMsg, _ := json.Marshal(
		Message{RPC: "JOIN", Msg: JoinMsg{Msg: sender.Me.String()}, Author: sender.Me.String()},
	)
	fmt.Println(string(joinMsg))
	conn.Write(joinMsg)

	var msg JoinRespMsg
	err = decoder.Decode(&msg)
	fmt.Printf("\nRESPONSE RECEIVED:\n\tRPC: JOIN\n\tSuccess: %v\n\tID: %v\n\tIP: %v",
		msg.Success,
		msg.ID,
		msg.IP,
	)
	if err != nil {
		ch <- false
		return
	}
	key, err := hashing.ToKademliaID(msg.ID)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	node := nodeutils.Node{ID: key, IP: msg.IP}
	sender.AddNode <- nodeutils.AddNodeOp{AddedNode: node}

	// Run LookUp on myself to fill up k-buckets according to Kademlia specification.
	sender.LookUp(sender.Me.ID)

	ch <- msg.Success
}
