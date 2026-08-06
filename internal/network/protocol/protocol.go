package protocol

import (
	"encoding/binary"
	"errors"
	"io"
	"message-broker/internal/config"
)

const (
	headerSize = 10
)

const (
	MagicZAKA uint32 = 0x5A414B41
)

const (
	TypePing uint8 = iota + 1
	TypePong
	TypeData
	TypeAck
	TypeError
)

var (
	ErrInvalidMagic           = errors.New("invalid magic bytes")
	ErrPayloadTooLarge        = errors.New("payload exceeds maximum size")
	ErrUnsupportedVersion     = errors.New("unsupported protocol version")
	ErrUnsupportedMessageType = errors.New("unsupported message type")
	ErrNoNeedPayload          = errors.New("no need payload for this type")
)

type ZProtocol struct {
	cfg            *config.ProtocolConfig
	maxPayloadSize uint32
	Version        uint8
}

func NewZProtocol(cfg *config.ProtocolConfig) *ZProtocol {
	if cfg.MaxPayLoadSize <= 0 {
		panic("payload size should be greater than 0")
	}
	return &ZProtocol{
		cfg:            cfg,
		maxPayloadSize: cfg.MaxPayLoadSize,
		Version:        cfg.Version,
	}
}

type Header struct {
	Version    uint8
	Type       uint8
	PayloadLen uint32
}

type Message struct {
	Header  Header
	Payload []byte
}

func (p *ZProtocol) Encode(w io.Writer, msg *Message) error {
	if len(msg.Payload) > int(p.maxPayloadSize) {
		return ErrPayloadTooLarge
	}

	header := make([]byte, headerSize)

	binary.BigEndian.PutUint32(header[0:4], MagicZAKA)

	header[4] = msg.Header.Version
	header[5] = msg.Header.Type

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

func (p *ZProtocol) DecodeHeader(r io.Reader) (*Header, error) {
	header := make([]byte, headerSize)

	if _, err := io.ReadFull(r, header); err != nil {
		return nil, err
	}

	magic := binary.BigEndian.Uint32(header[0:4])
	if magic != MagicZAKA {
		return nil, ErrInvalidMagic
	}

	version := header[4]
	if version != byte(p.Version) {
		return nil, ErrUnsupportedVersion
	}

	msgType := header[5]
	payloadLen := binary.BigEndian.Uint32(header[6:10])
	return &Header{
		Version:    version,
		Type:       msgType,
		PayloadLen: payloadLen,
	}, nil
}

func (p *ZProtocol) Decode(r io.Reader, header *Header) (*Message, error) {
	if header.PayloadLen > uint32(p.maxPayloadSize) {
		return nil, ErrPayloadTooLarge
	}

	payload := make([]byte, header.PayloadLen)
	if int(header.PayloadLen) > 0 {
		if _, err := io.ReadFull(r, payload); err != nil {
			return nil, err
		}
	}
	return &Message{
		Header:  *header,
		Payload: payload,
	}, nil
}

func isSupportedType(msgType byte) bool {
	supportedTypes := []uint8{TypeAck, TypePing, TypeData, TypeError, TypePong}
	for _, sType := range supportedTypes {
		if uint8(msgType) == sType {
			return true
		}
	}
	return false
}
