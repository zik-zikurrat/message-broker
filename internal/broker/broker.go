package broker

import (
	"log"
	"message-broker/internal/config"
)

type ZBroker struct {
	cfg *config.Config
}

func NewZBroker(cfg *config.Config) *ZBroker {
	return &ZBroker{
		cfg: cfg,
	}
}

func (b *ZBroker) Handle(payload []byte) error {
	log.Printf("broker received message")
	return nil
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

// type MessageInterface interface {
// 	Read([]byte) error
// 	Write() error
// 	Send() error
// }

// type TopicInterface interface {
// 	Purge() error
// 	Rebalance() error
// }

// type Message struct {
// 	Key     string
// 	Headers map[any]any
// 	Value   any
// 	offset  int64
// }

// type Partition struct {
// 	m        sync.RWMutex
// 	offset   int64
// 	Messages map[int64]Message
// }

// func NewPartition() *Partition {
// 	return &Partition{
// 		m:        sync.RWMutex{},
// 		offset:   0,
// 		Messages: make(map[int64]Message, 100),
// 	}
// }

// type ConsumerGroup struct {
// 	Consumers []struct{}
// 	Count     int64
// }

// type Topic struct {
// 	Name          string
// 	ConsumerGroup *ConsumerGroup
// 	Partitions    []*Partition
// }

// func NewTopic(name string) *Topic {
// 	return &Topic{
// 		Name:          name,
// 		ConsumerGroup: nil,
// 		Partitions:    nil,
// 	}
// }

// type Broker struct {
// 	Topics []*Topic
// }
