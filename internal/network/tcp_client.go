package network

import (
	"bufio"
	"fmt"
	"log"
	"message-broker/internal/network/protocol"
	"net"
	"time"
)

type TCPClient struct {
	connection   net.Conn
	protocol     *protocol.ZProtocol
	reader       *bufio.Reader
	writeTimeout time.Duration
	readTimeout  time.Duration
}

func NewTCPClinet(protocol *protocol.ZProtocol, address string, writeTimeout time.Duration, readTimeout time.Duration) (*TCPClient, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Printf("failed to create tcp client: %v\n", err)
		return nil, err
	}
	reader := bufio.NewReader(conn)
	return &TCPClient{
		connection:   conn,
		protocol:     protocol,
		reader:       reader,
		writeTimeout: writeTimeout,
		readTimeout:  readTimeout,
	}, nil
}

func (c *TCPClient) Send(msg *protocol.Message) (*protocol.Message, error) {
	if err := c.connection.SetWriteDeadline(time.Now().Add(c.writeTimeout)); err != nil {
		return nil, fmt.Errorf("set deadline: %w", err)
	}
	if err := c.protocol.Encode(c.connection, msg); err != nil {
		return nil, fmt.Errorf("encode message: %w", err)
	}
	if err := c.connection.SetReadDeadline(time.Now().Add(c.readTimeout)); err != nil {
		return nil, fmt.Errorf("set deadline: %w", err)
	}
	header, err := c.protocol.DecodeHeader(c.connection)
	if err != nil {
		return nil, fmt.Errorf("decode header: %w", err)
	}
	resp, err := c.protocol.Decode(c.connection, header)
	if err != nil {
		return nil, fmt.Errorf("docode: %w", err)
	}
	return resp, nil
}

func (c *TCPClient) Close() {
	if c.connection != nil {
		if err := c.connection.Close(); err != nil {
			log.Printf("client close: %v", err)
		}
	}
}
