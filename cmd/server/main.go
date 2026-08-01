package main

import (
	"context"
	"log"
	"message-broker/internal/config"
	"message-broker/internal/network"
	"message-broker/internal/network/protocol"
)

func main() {
	cfg := config.MustLoad()
	protocol := protocol.NewZProtocol(&cfg.Protocol)
	tcpServer, err := network.NewTCPServer(cfg, protocol)
	if err != nil {
		log.Printf("failed to create tcp server: %v", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tcpServer.Start(ctx)
}
