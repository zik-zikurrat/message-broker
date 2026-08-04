package main

import (
	"fmt"
	"log"
	"message-broker/internal/config"
	"message-broker/internal/network"
	"message-broker/internal/network/protocol"
)

func main() {
	cfg := config.MustLoad()
	zProtocol := protocol.NewZProtocol(&cfg.Protocol)
	tcpClient, err := network.NewTCPClinet(zProtocol, cfg.Server.Address, cfg.Connection.WriteTimeout, cfg.Connection.ReadTimeout)
	if err != nil {
		log.Printf("failed create client: %v", err)
		return
	}
	msg := protocol.Message{
		Header: protocol.Header{
			Version:    cfg.Protocol.Version,
			Type:       protocol.TypePing,
			PayloadLen: 0,
		},
		Payload: nil,
	}
	resp, err := tcpClient.Send(&msg)
	if err != nil {
		log.Printf("send error: %v\n", err)
	}
	fmt.Printf("RESP: %v\n", resp)
}
