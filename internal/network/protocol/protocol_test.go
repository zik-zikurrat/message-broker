package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"message-broker/internal/config"
	"testing"
)

const (
	InvalidMagic uint32 = 0x5A414B43
)

func InvalidEncode(w io.Writer, msg *Frame, p *ZProtocol) error {
	if len(msg.Payload) > int(p.maxPayloadSize) {
		return ErrPayloadTooLarge
	}

	header := make([]byte, headerSize)

	binary.BigEndian.PutUint32(header[0:4], InvalidMagic)

	header[4] = msg.FrameHeader.Version
	header[5] = msg.FrameHeader.Type

	binary.BigEndian.PutUint32(header[6:10], uint32(len(msg.Payload)))

	if _, err := w.Write(header); err != nil {
		return err
	}

	if len(msg.Payload) > 0 {
		if _, err := w.Write(msg.Payload); err != nil {
			return err
		}
	}
	return nil
}

func initializeProtocol() *ZProtocol {
	cfg := config.ProtocolConfig{
		HeartbeatSupported: false,
		MaxPayLoadSize:     1048576,
		Version:            1,
	}
	protocol := NewZProtocol(&cfg)
	return protocol
}

func TestDecodeHeaderMessageSuccess(t *testing.T) {
	var buf bytes.Buffer
	p := initializeProtocol()
	msg := Frame{
		FrameHeader: FrameHeader{
			Version:    1,
			Type:       TypePing,
			PayloadLen: 0,
		},
		Payload: nil,
	}
	p.Encode(&buf, &msg)
	_, err := p.DecodeHeader(&buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
func TestDecodeHeaderMessageFailureVersion(t *testing.T) {
	var buf bytes.Buffer
	p := initializeProtocol()
	msg := Frame{
		FrameHeader: FrameHeader{
			Version:    2,
			Type:       TypePing,
			PayloadLen: 0,
		},
		Payload: nil,
	}
	p.Encode(&buf, &msg)
	_, err := p.DecodeHeader(&buf)
	if err != nil {
		if errors.Is(err, ErrUnsupportedVersion) {
			return
		}
		t.Fatalf("unexpected error: expect ErrUnsupportedVersion, got: %v", err)
	}
}

func TestDecodeHeaderMessageFailureMagic(t *testing.T) {
	var buf bytes.Buffer
	p := initializeProtocol()
	msg := Frame{
		FrameHeader: FrameHeader{
			Version:    1,
			Type:       TypePing,
			PayloadLen: 0,
		},
		Payload: nil,
	}

	InvalidEncode(&buf, &msg, p)

	_, err := p.DecodeHeader(&buf)
	if err != nil {
		if errors.Is(err, ErrInvalidMagic) {
			return
		}
		t.Fatalf("unexpected error: expect ErrInvalidMagic, got: %v", err)
	}
}

func TestDecodeHeaderMessageFailurePayloadSize(t *testing.T) {
	var buf bytes.Buffer
	p := initializeProtocol()
	msg := Frame{
		FrameHeader: FrameHeader{
			Version:    1,
			Type:       TypePing,
			PayloadLen: 0,
		},
		Payload: nil,
	}

	p.Encode(&buf, &msg)

	_, err := p.DecodeHeader(&buf)
	if err != nil {
		t.Fatalf("unexpected error: expect ErrPayloadTooLarge, got: %v", err)
	}
	msg.FrameHeader.PayloadLen = p.maxPayloadSize + 1
	_, err = p.Decode(&buf, &msg.FrameHeader)
	if err != nil {
		if errors.Is(err, ErrPayloadTooLarge) {
			return
		}
		t.Fatalf("unexpected error: expect ErrPayloadTooLarge, got: %v", err)
	}
}
