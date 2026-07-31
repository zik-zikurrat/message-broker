package protocol

import (
	"encoding/binary"
	"errors"
	"io"
)

const (
	MagicByte1      = 0xCA
	MagicByte2      = 0xFE
	HeaderSize      = 8
	MaxPayLoadSize  = 1 << 20 // 1MB limit
	ProtocolVersion = 1
)

const (
	TypePing uint8 = iota + 1
	TypePong
	TypeData
	TypeAck
	TypeError
)

var (
	ErrInvalidMagic       = errors.New("invalid magic bytes")
	ErrPayloadTooLarge    = errors.New("payload exceeds maximum size")
	ErrUnsupportedVersion = errors.New("unsupported protocol version")
)

type Message struct {
	Version uint8
	Type    uint8
	Payload []byte
}

func Encode(w io.Writer, msg *Message) error {
	if len(msg.Payload) > MaxPayLoadSize {
		return ErrPayloadTooLarge
	}

	header := make([]byte, HeaderSize)

	header[0] = MagicByte1
	header[1] = MagicByte2

	header[2] = msg.Version
	header[3] = msg.Type

	binary.BigEndian.PutUint32(header[4:8], uint32(len(msg.Payload)))

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

func Decode(r io.Reader) (*Message, error) {
	header := make([]byte, HeaderSize)

	if _, err := io.ReadFull(r, header); err != nil {
		return nil, err
	}

	if header[0] != MagicByte1 || header[1] != MagicByte2 {
		return nil, ErrInvalidMagic
	}

	version := header[2]
	if version != ProtocolVersion {
		return nil, ErrUnsupportedVersion
	}

	msgType := header[3]
	payloadLen := binary.BigEndian.Uint32(header[4:8])

	if payloadLen > MaxPayLoadSize {
		return nil, ErrPayloadTooLarge
	}

	payload := make([]byte, payloadLen)
	if payloadLen > 0 {
		if _, err := io.ReadFull(r, payload); err != nil {
			return nil, err
		}
	}

	return &Message{
		Version: version,
		Type:    msgType,
		Payload: payload,
	}, nil
}
