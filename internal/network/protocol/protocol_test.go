package protocol

// import (
// 	"bytes"
// 	"errors"
// 	"message-broker/internal/config"
// 	"testing"
// )

// func initializeProtocol() *ZProtocol {
// 	cfg := config.ProtocolConfig{
// 		HeartbeatSupported: false,
// 		MaxPayLoadSize:     1048576,
// 		Version:            1,
// 	}
// 	protocol := NewZProtocol(&cfg)
// 	return protocol
// }

// func TestMessageSendSuccess(t *testing.T) {
// 	var buf bytes.Buffer
// 	p := initializeProtocol()
// 	p.Encode(&buf, &Message{Version: 1, Type: TypePing})
// 	_, err := p.Decode(&buf)
// 	if err != nil {
// 		t.Fatalf("unexpected error: %v", err)
// 	}
// }

// func TestMessageSendFailed(t *testing.T) {
// 	var buf bytes.Buffer
// 	p := initializeProtocol()
// 	p.Encode(&buf, &Message{Version: 2, Type: TypePing})
// 	_, err := p.Decode(&buf)
// 	if err != nil {
// 		if errors.Is(err, ErrUnsupportedVersion) {
// 			return
// 		} else {
// 			t.Fatalf("unexpected error: expect ErrUnsupportedVersion, got: %v", err)
// 		}
// 	}
// }
// func TestEncodeMessageWithPayloadSuccess(t *testing.T) {
// 	var buf bytes.Buffer
// 	p := initializeProtocol()
// 	if err := p.Encode(&buf, &Message{Version: 1, Type: TypeData, Payload: []byte("Test message")}); err != nil {
// 		t.Fatalf("unexpected error: %v", err)
// 	}
// }

// func TestDecodeMessageWithPayloadSuccess(t *testing.T) {
// 	var buf bytes.Buffer
// 	p := initializeProtocol()
// 	payload := bytes.Repeat([]byte{'X'}, 1048576)
// 	if err := p.Encode(&buf, &Message{Version: 1, Type: TypeData, Payload: payload}); err != nil {
// 		t.Fatalf("unexpected error: %v", err)
// 	}
// 	_, err := p.Decode(&buf)
// 	if err != nil {
// 		t.Fatalf("unexpected error: %v", err)
// 	}
// }

// func TestEncodeMessageWithPayloadFailed(t *testing.T) {
// 	var buf bytes.Buffer
// 	p := initializeProtocol()
// 	payload := bytes.Repeat([]byte{'X'}, 1048577)
// 	if err := p.Encode(&buf, &Message{Version: 1, Type: TypeData, Payload: payload}); err != nil {
// 		if errors.Is(err, ErrPayloadTooLarge) {
// 			return
// 		} else {
// 			t.Fatalf("unexpected error: expect ErrPayloadTooLarge, got: %v", err)
// 		}
// 	}
// }
