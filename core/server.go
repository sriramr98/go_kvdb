package core

import (
	"errors"
	"fmt"
	"io"
	"net"
)

type ServerOpts struct {
	Port int
}

type Server struct {
	opts ServerOpts
}

func NewServer(opts ServerOpts) *Server {
	return &Server{opts: opts}
}

func (s *Server) Start() {
	fmt.Println("Starting server on port", s.opts.Port)
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", s.opts.Port))
	if err != nil {
		fmt.Printf("Error %s", err)
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
			fmt.Printf("Error %s", err)
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
		fmt.Printf("Received %d bytes: %s\n", n, buf[:n])
	}
}
