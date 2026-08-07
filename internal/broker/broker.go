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

type Message struct {
	Key     []byte
	Value   []byte
	Headers map[any]any
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
