package producer

import "message-broker/internal/network/protocol"

type ProducerRequest struct {
	Topic     string
	Partition *int32
	Key       []byte
	Value     []byte
	Headers   map[any]any
}

func ToProducerRequest(frame *protocol.Frame) *ProducerRequest {

}
