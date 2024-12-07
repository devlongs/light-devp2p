package node

import (
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/p2p"
)

type Peer struct {
	*p2p.Peer
	rw      p2p.MsgReadWriter
	trusted bool
	quit    chan struct{}
}

type PeerManager struct {
	peers   map[string]*Peer
	lock    sync.RWMutex
	metrics *Metrics
}

func NewPeer(p *p2p.Peer, rw p2p.MsgReadWriter) *Peer {
	return &Peer{
		Peer: p,
		rw:   rw,
		quit: make(chan struct{}),
	}
}

func NewPeerManager() *PeerManager {
	return &PeerManager{
		peers: make(map[string]*Peer),
	}
}

func (pm *PeerManager) Handle(p *Peer) error {
	pm.lock.Lock()
	pm.peers[p.ID().String()] = p
	pm.lock.Unlock()

	// Send hello message
	if err := p2p.Send(p.rw, HelloMsg, &HelloMessage{
		Version: 1,
		NodeID:  p.ID().String(),
	}); err != nil {
		return fmt.Errorf("failed to send hello: %v", err)
	}

	defer func() {
		pm.lock.Lock()
		delete(pm.peers, p.ID().String())
		pm.lock.Unlock()
	}()

	return pm.readLoop(p)
}

func (pm *PeerManager) readLoop(p *Peer) error {
	for {
		msg, err := p.rw.ReadMsg()
		if err != nil {
			return fmt.Errorf("failed to read message: %v", err)
		}

		switch msg.Code {
		case HelloMsg:
			var hello HelloMessage
			if err := msg.Decode(&hello); err != nil {
				return fmt.Errorf("failed to decode hello: %v", err)
			}
			fmt.Printf("Received hello from %s (version %d)\n", hello.NodeID, hello.Version)
		default:
			fmt.Printf("Received unknown message %d from %s\n", msg.Code, p.ID())
		}

		if pm.metrics != nil {
			pm.metrics.Lock()
			pm.metrics.messagesReceived++
			pm.metrics.Unlock()
		}

		msg.Discard()
	}
}

func (pm *PeerManager) Broadcast(msgCode uint64, data interface{}) error {
	pm.lock.RLock()
	defer pm.lock.RUnlock()

	for _, peer := range pm.peers {
		if err := p2p.Send(peer.rw, msgCode, data); err != nil {
			fmt.Printf("Failed to send message to peer %s: %v\n", peer.ID(), err)
			continue
		}

		if pm.metrics != nil {
			pm.metrics.Lock()
			pm.metrics.messagesSent++
			pm.metrics.Unlock()
		}
	}

	return nil
}
