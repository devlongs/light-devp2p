package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/devlongs/light-devp2p/config"
	"github.com/devlongs/light-devp2p/node"
)

func main() {
	cfg := config.ParseFlags()

	// Create and start the node
	n, err := node.NewNode(cfg)
	if err != nil {
		log.Fatalf("Failed to create node: %v", err)
	}

	if err := n.Start(); err != nil {
		log.Fatalf("Failed to start node: %v", err)
	}
	defer n.Stop()

	// Connect to bootnode if specified
	if cfg.BootNodeURL != "" {
		if err := n.ConnectToBootNode(cfg.BootNodeURL); err != nil {
			log.Printf("Failed to connect to bootnode: %v", err)
		}
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Println("Shutting down...")
}
