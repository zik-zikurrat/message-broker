package producer

type ProducerRequest struct {
	Topic     string
	Partition *int32
	Key       []byte
	Value     []byte
	Headers   map[any]any
}
