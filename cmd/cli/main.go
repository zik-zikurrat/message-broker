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
	payload := []byte("Hello world!")
	msg := protocol.Frame{
		Header: protocol.FrameHeader{
			Version:    1,
			Type:       protocol.TypeProduce,
			PayloadLen: uint32(len(payload)),
		},
		Payload: payload,
	}
	resp, err := tcpClient.Send(&msg)
	if err != nil {
		log.Printf("send error: %v\n", err)
	}
	fmt.Printf("RESP: Header: %v, Payload: %b\n", resp.Header, resp.Payload)
}
