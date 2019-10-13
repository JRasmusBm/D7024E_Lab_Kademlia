package network

type Message struct {
	RPC string
	Author string
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

type JoinRespMsg struct {
	Success bool
	ID string
}

// FIND_VALUE message
type FindValueMsg struct {
	Key string
}

type FindValueRespMsg struct {
	Content string
	Nodes string
}
