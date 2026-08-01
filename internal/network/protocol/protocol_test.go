package protocol

import (
	"bytes"
	"message-broker/internal/config"
	"testing"
)

func initializeProtocol() *ZProtocol {
	cfg := config.ProtocolConfig{
		HeartbeatSupported: false,
		MaxPayLoadSize:     1048576,
		Version:            1,
	}
	protocol := NewZProtocol(&cfg)
	return protocol
}

func TestMessageTypePing(t *testing.T) {
	var buf bytes.Buffer
	p := initializeProtocol()
	p.Encode(&buf, &Message{Version: 1, Type: TypePing})
	got, err := p.Decode(&buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Type != TypePong {
		t.Fatal("unexpeted response type")
	}
}

func TestDecodeMessage(t *testing.T) {

}
