## light-devp2p
A lightweight implementation of Ethereum's devp2p protocol. This project provides a minimal, educational implementation of peer-to-peer networking using the core concepts of Ethereum's network layer.

## Features
- Minimal devp2p implementation
- Node discovery and peer management
- Secure communication using ECDSA keys
- Configurable network parameters
- debug logging
- Bootnode supports

### Installation
```bash
git clone https://github.com/devlongs/light-devp2p.git
cd light-devp2p
go mod tidy
```

### Running Nodes
1. Start the first node:
```bash
go run main.go -addr :30303
```

2. Start the second node:
```bash
go run main.go -addr :30304 -bootnode "enode://[public-key]@127.0.0.1:30303"
```

### Configuration
Available command-line options:

| Flag       | Description                 | Default          |
|------------|-----------------------------|------------------|
| `-addr`    | Listen address              | `":30303"`       |
| `-bootnode`| Bootnode enode URL          | `""`             |
| `-maxpeers`| Maximum peer connections    | `10`             |
| `-nodename`| Node identifier             | `"MinimalP2PNode"` |
| `-metrics` | Enable metrics collection   | `false`          |

### Project Structure
```txt
light-devp2p/
├── main.go           # Entry point
├── go.mod           # Go module file
├── node/
│   ├── node.go      # Core node implementation
│   ├── peer.go      # Peer management
│   └── protocol.go  # Protocol implementation
└── config/
    └── config.go    # Configuration handling
```

### Main Components
The implementation is organized around three main components:

1. Node

- Manages the P2P server
- Handles peer connections
- Implements protocol communication

2. Peer Manager

- Tracks active connections
- Manages peer lifecycle
- Handles peer state

3. Protocol

- Defines message types
- Implements handshake
- Manages message exchange
