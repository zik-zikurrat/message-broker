package main

import (
	"io"
	"log"
	"message-broker/internal/protocol"
	"net"
	"sync"
)

type MessageInterface interface {
	Read([]byte) error
	Write() error
	Send() error
}

type TopicInterface interface {
	Purge() error
	Rebalance() error
}

type Message struct {
	Key     string
	Headers map[any]any
	Value   any
	offset  int64
}

type Partition struct {
	m        sync.RWMutex
	offset   int64
	Messages map[int64]Message
}

func NewPartition() *Partition {
	return &Partition{
		m:        sync.RWMutex{},
		offset:   0,
		Messages: make(map[int64]Message, 100),
	}
}

type ConsumerGroup struct {
	Consumers []struct{}
	Count     int64
}

type Topic struct {
	Name          string
	ConsumerGroup *ConsumerGroup
	Partitions    []*Partition
}

func NewTopic(name string) *Topic {
	return &Topic{
		Name:          name,
		ConsumerGroup: nil,
		Partitions:    nil,
	}
}

type Broker struct {
	Topics            []*Topic
	ISR               int64
	ReplicationFactor int64
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		msg, err := protocol.Decode(conn)
		if err != nil {
			if err != io.EOF {
				log.Printf("decode error: %v\n", err)
			}
			return
		}
		switch msg.Type {
		case protocol.TypePing:
			response := &protocol.Message{
				Version: protocol.ProtocolVersion,
				Type:    protocol.TypePong,
				Payload: nil,
			}
			if err := protocol.Encode(conn, response); err != nil {
				log.Printf("encode error: %v\n", err)
				return
			}
		case protocol.TypeData:
			processData(msg.Payload)

			ack := &protocol.Message{
				Version: protocol.ProtocolVersion,
				Type:    protocol.TypeAck,
				Payload: nil,
			}
			if err := protocol.Encode(conn, ack); err != nil {
				log.Printf("encode error: %v\n", err)
				return
			}
		}
	}
}

func processData(payload []byte) {
	log.Printf("received %d bytes", len(payload))
}

func main() {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Println("Server listening on :9000")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept error: %v\n", err)
			continue
		}
		go handleConnection(conn)
	}
}

// у сообщения есть key, headers, value
// сообщение имеет свой личный номер
// сообщение можно читать, отправить, записать в лог
// сообщение можно отправит в топик
// у топика есть много партиций, партиции читают консумеры
// топик должен знать сколько consumers его читают

// брокер владеет топиками и управляет репликами
// топик владеет партициями и управляет их колличеством
// партиции владеют сообщениями и управляют доступом к ним
