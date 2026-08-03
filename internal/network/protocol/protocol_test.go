package protocol

import (
	"bytes"
	"errors"
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

func TestMessageSendSuccess(t *testing.T) {
	var buf bytes.Buffer
	p := initializeProtocol()
	p.Encode(&buf, &Message{Version: 1, Type: TypePing})
	_, err := p.Decode(&buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMessageSendFailedVersion(t *testing.T) {
	var buf bytes.Buffer
	p := initializeProtocol()
	p.Encode(&buf, &Message{Version: 2, Type: TypePing})
	_, err := p.Decode(&buf)
	if err != nil {
		if errors.Is(err, ErrUnsupportedVersion) {
			return
		} else {
			t.Fatalf("unexpected error: expect ErrUnsupportedVersion, got: %v", err)
		}
	}
}
func TestDecodeMessage(t *testing.T) {

}
