package core

import (
	"errors"
	"fmt"
	"io"
	"net"

	"gitub.com/sriramr98/go_kvdb/core/processors"
	"gitub.com/sriramr98/go_kvdb/core/protocol"
)

type ServerOpts struct {
	Port int
}

type Server struct {
	opts             ServerOpts
	requestProcessor processors.RequestProcessor
	protocol         protocol.Protocol
}

func NewServer(opts ServerOpts, processor processors.RequestProcessor, protocol protocol.Protocol) *Server {
	return &Server{opts: opts, requestProcessor: processor, protocol: protocol}
}

func (s *Server) Start() {
	fmt.Println("Starting server on port", s.opts.Port)
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", s.opts.Port))
	if err != nil {
		fmt.Printf("Error %s\n", err)
		return
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("Connection closed")
				continue
			}
			fmt.Printf("Error %s\n", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("New connection")
	for {
		buf := make([]byte, 2048)
		n, err := conn.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("Connection closed")
				return
			}
			fmt.Printf("Error %s", err)
			return
		}

		data_received := string(buf[:n])
		// The last character is a newline, so we remove it
		data_received = data_received[:len(data_received)-1]
		fmt.Printf("Received %d bytes: %s\n", n, data_received)

		request, err := s.protocol.Parse(data_received)
		if err != nil {
			fmt.Printf("Error parsing protocol %s\n", err)
			s.writeError(err, conn)
			continue
		}
		response, err := s.requestProcessor.Process(request)
		if err != nil {
			fmt.Printf("Error processing request %s\n", err)
			s.writeError(err, conn)
			continue
		}

		s.writeSuccess(response.Value, conn)
	}
}

func (s *Server) writeError(err error, conn net.Conn) {
	conn.Write([]byte(fmt.Sprintf("ERR: %s\n", err)))
}

func (s *Server) writeSuccess(value []byte, conn net.Conn) {
	if len(value) == 0 {
		conn.Write([]byte("OK\n"))
		return
	}
	conn.Write([]byte(fmt.Sprintf("OK: %s\n", value)))
}
