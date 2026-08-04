package network

import (
	"context"
	"errors"
	"log"
	"message-broker/internal/config"
	"message-broker/internal/network/protocol"
	"net"
	"sync"
	"time"
)

type TCPServer struct {
	listener       net.Listener
	idleTimeout    time.Duration
	readTimeout    time.Duration
	writeTimeout   time.Duration
	maxConnections int
	protocol       *protocol.ZProtocol
}

func NewTCPServer(cfg *config.Config, protocol *protocol.ZProtocol) (*TCPServer, error) {
	listener, err := net.Listen("tcp", cfg.Server.Address)
	if err != nil {
		return nil, err
	}
	return &TCPServer{
		listener:       listener,
		idleTimeout:    cfg.Connection.IdleTimeout,
		readTimeout:    cfg.Connection.ReadTimeout,
		writeTimeout:   cfg.Connection.WriteTimeout,
		maxConnections: cfg.Server.MaxConnections,
		protocol:       protocol,
	}, nil
}

func (s *TCPServer) Start(ctx context.Context) {
	sem := make(chan struct{}, s.maxConnections)
	wg := sync.WaitGroup{}
	log.Printf("server start with address: %s\n", s.listener.Addr())

	go func() {
		<-ctx.Done()
		if err := s.listener.Close(); err != nil {
			log.Printf("close server: %v\n", err)
		}
	}()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if ctx.Err() != nil {
				wg.Wait()
				return
			}
			log.Printf("error while accepting connection: %v", err)
			continue
		}
		select {
		case sem <- struct{}{}:
		case <-ctx.Done():
			if err := conn.Close(); err != nil {
				log.Printf("close server: %v\n", err)
			}
			wg.Wait()
			return
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			s.handleConnection(conn)
		}()
	}
}

func (s *TCPServer) handleConnection(conn net.Conn) {
	defer func() {
		if v := recover(); v != nil {
			log.Printf("captured panic: %v\n", v)
		}
	}()
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("close server: %v\n", err)
		}
	}()
	for {
		conn.SetReadDeadline(time.Now().Add(s.idleTimeout))
		header, err := s.protocol.DecodeHeader(conn)
		if err != nil {
			s.sendError(conn, err)
			break
		}
		conn.SetReadDeadline(time.Now().Add(s.readTimeout))
		msg, err := s.protocol.Decode(conn, header)
		if err != nil {
			s.sendError(conn, err)
			break
		}
		conn.SetWriteDeadline(time.Now().Add(s.writeTimeout))
		if err := s.protocol.Encode(conn, msg); err != nil {
			s.sendError(conn, err)
			break
		}
	}
}

func isSendbleError(err error) bool {
	protocolErrors := []error{protocol.ErrInvalidMagic, protocol.ErrPayloadTooLarge, protocol.ErrUnsupportedVersion}
	for _, pErr := range protocolErrors {
		if errors.Is(err, pErr) {
			return true
		}
	}
	return false
}

func (s *TCPServer) sendError(conn net.Conn, err error) {
	if !isSendbleError(err) {
		return
	}
	log.Printf("got protocol error: %v\n", err)
	errMsg := &protocol.Message{
		Header: protocol.Header{
			Version:    s.protocol.Version,
			Type:       protocol.TypeError,
			PayloadLen: uint32(len([]byte(err.Error()))),
		},
		Payload: []byte(err.Error()),
	}
	if err := conn.SetWriteDeadline(time.Now().Add(s.writeTimeout)); err != nil {
		log.Printf("failed to set write deadline: %v", err)
	}
	if err := s.protocol.Encode(conn, errMsg); err != nil {
		log.Printf("failed to send error: %v", err)
	}
}
