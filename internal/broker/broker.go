package broker

import (
	"log"
	"message-broker/internal/config"
	"message-broker/internal/network/protocol"
)

type ZBroker struct {
	cfg *config.Config
}

func NewZBroker(cfg *config.Config) *ZBroker {
	return &ZBroker{
		cfg: cfg,
	}
}

func (b *ZBroker) Handle(frame *protocol.Frame) error {
	log.Printf("broker received message")
	switch frame.FrameHeader.Type {
	case protocol.TypeFetch:
		panic("not implemented logic for fetch type")
	case protocol.TypeProduce:
		panic("not implemented logic for produce type")
	}
	return nil
}

type Message struct {
	Key     []byte
	Value   []byte
	Headers map[any]any
}
