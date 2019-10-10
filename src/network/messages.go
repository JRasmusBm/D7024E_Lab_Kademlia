package network

type Message struct {
	RPC string
	Msg interface{}
}

type AckMsg struct {
	Success bool
}

// PING message
type PingMsg struct {
	Msg string
}

// FIND_NODE messages
type FindNodeMsg struct {
	ID string
}
type FindNodeRespMsg struct {
	Nodes string
}

// Store message
type StoreMsg struct {
	Data string
}

// JOIN message
type JoinMsg struct {
	Msg string
}
