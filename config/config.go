package config

import (
	"flag"
	"os"
	"path/filepath"
)

type Config struct {
	ListenAddr    string
	BootNodeURL   string
	NodeKeyFile   string
	MaxPeers      int
	ProtocolName  string
	NodeName      string
	NetworkID     uint64
	EnableMetrics bool
}

func ParseFlags() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.ListenAddr, "addr", ":30303", "listen address")
	flag.StringVar(&cfg.BootNodeURL, "bootnode", "", "bootnode enode URL")
	flag.StringVar(&cfg.NodeKeyFile, "nodekey", filepath.Join(os.TempDir(), "node.key"), "node key file")
	flag.IntVar(&cfg.MaxPeers, "maxpeers", 10, "maximum number of peers")
	flag.StringVar(&cfg.ProtocolName, "protocol", "minp2p", "protocol name")
	flag.StringVar(&cfg.NodeName, "nodename", "MinimalP2PNode", "node name")
	flag.Uint64Var(&cfg.NetworkID, "networkid", 1, "network ID")
	flag.BoolVar(&cfg.EnableMetrics, "metrics", false, "enable metrics collection")

	flag.Parse()
	return cfg
}
