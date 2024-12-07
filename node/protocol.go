package node

import (
	"time"
)

// Protocol message codes
const (
	HelloMsg = iota
	StatusMsg
	MessageMsg
)

// Protocol message types
type HelloMessage struct {
	Version uint64
	NodeID  string
	Time    uint64
}

type StatusMessage struct {
	NetworkID uint64
	Version   uint64
}

type Message struct {
	From    string
	Content string
	Time    uint64
}

func NewHelloMessage(nodeID string) *HelloMessage {
	return &HelloMessage{
		Version: 1,
		NodeID:  nodeID,
		Time:    uint64(time.Now().Unix()),
	}
}

func NewStatusMessage(networkID uint64) *StatusMessage {
	return &StatusMessage{
		NetworkID: networkID,
		Version:   1,
	}
}

func NewMessage(from, content string) *Message {
	return &Message{
		From:    from,
		Content: content,
		Time:    uint64(time.Now().Unix()),
	}
}
