package protocol

import (
	"encoding/binary"
	"errors"
	"io"
	"log"
	"message-broker/internal/config"
)

const (
	minHeaderSize = 8
)

const (
	MagicZA = 0x5A41
	MagicKA = 0x4B41
)

const (
	TypePing uint8 = iota + 1
	TypePong
	TypeData
	TypeAck
	TypeError
)

// Я должен не делать handshake просто слать сообщение брокер отвечает в зависимости от гарантии 0, 1, all
// Обязанность протокола - кодировать декодировать сообщение, валидировать его размер и прочие детали, так же, если heartbeat перестал отправляться запустить ребаласировку брокера

var (
	ErrInvalidMagic       = errors.New("invalid magic bytes")
	ErrPayloadTooLarge    = errors.New("payload exceeds maximum size")
	ErrUnsupportedVersion = errors.New("unsupported protocol version")
)

type ZProtocol struct {
	cfg            *config.ProtocolConfig
	maxPayloadSize uint32
	Version        uint8
	headerSize     uint8
}

func NewZProtocol(cfg *config.ProtocolConfig) *ZProtocol {
	if cfg.HeaderSize < minHeaderSize {
		panic("header size must be at least 8 bytes")
	}
	return &ZProtocol{
		cfg:            cfg,
		maxPayloadSize: cfg.MaxPayLoadSize,
		Version:        cfg.Version,
		headerSize:     cfg.HeaderSize,
	}
}

type Message struct {
	Version uint8
	Type    uint8
	Payload []byte
}

func (p *ZProtocol) Encode(w io.Writer, msg *Message) error {
	if len(msg.Payload) > int(p.maxPayloadSize) {
		return ErrPayloadTooLarge
	}

	header := make([]byte, p.headerSize)

	binary.BigEndian.PutUint16(header[0:2], MagicZA)

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

func (p *ZProtocol) Decode(r io.Reader) (*Message, error) {
	header := make([]byte, p.headerSize)

	if _, err := io.ReadFull(r, header); err != nil {
		return nil, err
	}

	magic := binary.BigEndian.Uint16(header[0:2])

	switch magic {
	case MagicZA:
		log.Println("header identified: type ZA")
	case MagicKA:
		log.Println("header identified: type KA")
	default:
		return nil, ErrInvalidMagic
	}

	version := header[2]
	if version != byte(p.Version) {
		return nil, ErrUnsupportedVersion
	}

	msgType := header[3]
	payloadLen := binary.BigEndian.Uint32(header[4:8])

	if payloadLen > p.maxPayloadSize {
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

func processData(payload []byte) {
	log.Printf("received %d bytes", len(payload))
}
