package node

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/devlongs/light-devp2p/config"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
)

type Node struct {
	cfg         *config.Config
	server      *p2p.Server
	privateKey  *ecdsa.PrivateKey
	peerManager *PeerManager
	protocols   []p2p.Protocol
	quit        chan struct{}
	wg          sync.WaitGroup
	metrics     *Metrics
}

type Metrics struct {
	sync.RWMutex
	messagesSent     uint64
	messagesReceived uint64
	peersConnected   uint64
	lastUpdate       time.Time
}

func NewNode(cfg *config.Config) (*Node, error) {
	// Load or generate private key
	privateKey, err := loadOrGenerateNodeKey(cfg.NodeKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load/generate node key: %v", err)
	}

	peerManager := NewPeerManager()

	// Create protocols
	protocols := []p2p.Protocol{
		{
			Name:    cfg.ProtocolName,
			Version: 1,
			Length:  1,
			Run: func(p *p2p.Peer, rw p2p.MsgReadWriter) error {
				peer := NewPeer(p, rw)
				return peerManager.Handle(peer)
			},
		},
	}

	// Create metrics if enabled
	var metrics *Metrics
	if cfg.EnableMetrics {
		metrics = &Metrics{lastUpdate: time.Now()}
	}

	// Create P2P server configuration
	serverConfig := p2p.Config{
		PrivateKey: privateKey,
		MaxPeers:   cfg.MaxPeers,
		Name:       cfg.NodeName,
		ListenAddr: cfg.ListenAddr,
		Protocols:  protocols,
	}

	return &Node{
		cfg:         cfg,
		server:      &p2p.Server{Config: serverConfig},
		privateKey:  privateKey,
		peerManager: peerManager,
		protocols:   protocols,
		quit:        make(chan struct{}),
		metrics:     metrics,
	}, nil
}

func (n *Node) Start() error {
	if err := n.server.Start(); err != nil {
		return fmt.Errorf("failed to start p2p server: %v", err)
	}

	nodeURL := n.getNodeURL()
	nodeID := n.server.Self().ID()
	fmt.Printf("Node started with ID: %s\n", nodeID.String())
	fmt.Printf("Enode URL: %s\n", nodeURL)

	// Start metrics collection if enabled
	if n.metrics != nil {
		n.wg.Add(1)
		go n.collectMetrics()
	}

	return nil
}

func (n *Node) getNodeURL() string {
	pubkey := n.privateKey.Public().(*ecdsa.PublicKey)

	// Convert to the correct format (64-byte public key without compression prefix)
	pubkeyBytes := crypto.FromECDSAPub(pubkey)
	if len(pubkeyBytes) == 65 && pubkeyBytes[0] == 4 {
		pubkeyBytes = pubkeyBytes[1:] // Remove the compression prefix
	}

	// Format the enode URL with the hex-encoded public key
	_, portStr, _ := net.SplitHostPort(n.cfg.ListenAddr)
	port := uint16(30303)
	if p, err := net.LookupPort("tcp", portStr); err == nil {
		port = uint16(p)
	}

	return fmt.Sprintf("enode://%s@127.0.0.1:%d", hex.EncodeToString(pubkeyBytes), port)
}

func (n *Node) Stop() error {
	close(n.quit)
	n.server.Stop()
	n.wg.Wait()
	return nil
}

func (n *Node) ConnectToBootNode(bootNodeURL string) error {
	node, err := enode.Parse(enode.ValidSchemes, bootNodeURL)
	if err != nil {
		return fmt.Errorf("failed to parse bootnode URL: %v", err)
	}
	n.server.AddPeer(node)
	return nil
}

func (n *Node) Broadcast(msgCode uint64, data interface{}) error {
	return n.peerManager.Broadcast(msgCode, data)
}

func (n *Node) collectMetrics() {
	defer n.wg.Done()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			n.metrics.Lock()
			fmt.Printf("Metrics - Messages Sent: %d, Received: %d, Peers: %d\n",
				n.metrics.messagesSent,
				n.metrics.messagesReceived,
				n.metrics.peersConnected,
			)
			n.metrics.Unlock()
		case <-n.quit:
			return
		}
	}
}

func loadOrGenerateNodeKey(keyfile string) (*ecdsa.PrivateKey, error) {
	if keyfile == "" {
		return crypto.GenerateKey()
	}

	if key, err := crypto.LoadECDSA(keyfile); err == nil {
		return key, nil
	}

	key, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	if err := crypto.SaveECDSA(keyfile, key); err != nil {
		return nil, err
	}

	return key, nil
}
