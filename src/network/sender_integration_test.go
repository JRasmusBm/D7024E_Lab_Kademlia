package network

// import (
// 	"testing"
// 	"utils/hashing"
// 	nodeutils "utils/node"
// )

// func TestDialIntegration(t *testing.T) {
// 	var dialer Dialer = &RealDialer{}
// 	var sender Sender = &RealSender{Dialer: &dialer}
// 	_, err := sender.Dial(nodeutils.Node{IP: "0.0.0.0"})
// 	if err == nil {
// 		t.Errorf("Should be unable to connect")
// 	}
// }

// func TestLookUpIntegration(t *testing.T) {
// 	findClosestNodes := make(chan nodeutils.FindClosestNodesOp, 1000)
// 	go func() {
// 		message := <-findClosestNodes
// 		message.Resp <- []nodeutils.Node{}
// 	}()
// 	var sender Sender = &RealSender{
// 		FindClosestNodes: findClosestNodes,
// 	}
// 	key, _ := hashing.ToKademliaID("0000000000000000000000000000000000000000")
// 	nodes := sender.LookUp(key)
// 	if len(nodes) > 0 {
// 		t.Errorf("Should not find any node")
// 	}
// }

// func TestLookUpIntegrationWithNodes(t *testing.T) {
// 	findClosestNodes := make(chan nodeutils.FindClosestNodesOp, 1000)
// 	id1, _ := hashing.ToKademliaID("1111111111111111111111111111111111111111")
// 	id2, _ := hashing.ToKademliaID("2222222222222222222222222222222222222222")
// 	id3, _ := hashing.ToKademliaID("3333333333333333333333333333333333333333")
// 	node1 := nodeutils.Node{ID: id1, IP: "0.0.0.0"}
// 	node2 := nodeutils.Node{ID: id2, IP: "0.0.0.0"}
// 	node3 := nodeutils.Node{ID: id3, IP: "0.0.0.0"}
// 	askNodes := []nodeutils.Node{
// 		node1,
// 		node2,
// 		node3,
// 	}
// 	go func() {
// 		message := <-findClosestNodes
// 		message.Resp <- askNodes
// 	}()
// 	var dialer Dialer = &RealDialer{}
// 	var sender Sender = &RealSender{
// 		FindClosestNodes: findClosestNodes,
// 		Dialer:           &dialer,
// 	}
// 	key, _ := hashing.ToKademliaID("0000000000000000000000000000000000000000")
// 	nodes := sender.LookUp(key)
//   if len(nodes) > 0 {
//     t.Errorf("Should not find any node")
//   }
// }
