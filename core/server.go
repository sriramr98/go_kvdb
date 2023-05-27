package core

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"

	"gitub.com/sriramr98/go_kvdb/core/processors"
	"gitub.com/sriramr98/go_kvdb/core/protocol"
	"gitub.com/sriramr98/go_kvdb/store"
)

type ServerRole int

const (
	ClientServerRole  ServerRole = iota
	ReplicaServerRole ServerRole = iota
)

type ServerOpts struct {
	Port       int
	Role       ServerRole
	IsLeader   bool
	LeaderAddr string
}

type Server struct {
	opts             ServerOpts
	requestProcessor processors.RequestProcessor
	protocol         protocol.Protocol
	followerStore    store.DataStorer[net.Conn, struct{}]
	leaderConn       net.Conn
}

func NewServer(opts ServerOpts, clientProcessor processors.RequestProcessor, followerStore store.DataStorer[net.Conn, struct{}], protocol protocol.Protocol, ctx context.Context) *Server {
	server := &Server{opts: opts, requestProcessor: clientProcessor, protocol: protocol, followerStore: followerStore}

	// If this is a follower, we don't wanna start the follower without making a connection to the leader
	if !opts.IsLeader {
		leaderConn, err := net.Dial("tcp", opts.LeaderAddr)
		if err != nil {
			fmt.Printf("Error connecting to leader %s\n", err)
			return nil
		}

		server.leaderConn = leaderConn
	}

	return server
}

func (s *Server) Start() {
	fmt.Println("Starting server on port", s.opts.Port)
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", s.opts.Port))
	if err != nil {
		fmt.Printf("Error %s\n", err)
		return
	}
	defer ln.Close()

	if !s.opts.IsLeader {
		s.syncWithLeader()
		go s.handleConnection(s.leaderConn)
	}

	s.listenForConnections(ln)
}

func (s *Server) listenForConnections(ln net.Listener) {
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

		if s.opts.Role == ReplicaServerRole {
			s.followerStore.Set(conn, struct{}{})
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) syncWithLeader() {
	_, err := s.leaderConn.Write([]byte("SYNC\n"))
	if err != nil {
		fmt.Printf("Error syncing with leader %s\n", err)
		return
	}

	buf := make([]byte, 2048)
	n, err := s.leaderConn.Read(buf)
	if err != nil {
		if errors.Is(err, io.EOF) {
			fmt.Println("Connection closed")
			return
		}
		fmt.Printf("Error %s", err)
		return
	}

	data_received := buf[:n]
	fmt.Printf("Received state %s\n", string(data_received))

	//TODO: Deserialize state

	//TODO: Set state in client store
	// s.requestProcessor.SetAll()
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

		if s.opts.Role == ClientServerRole && request.Command.IsReplicable {
			go s.propogateToFollowers(data_received)
		}

		s.writeSuccess(response.Value, conn)
	}
}

func (s *Server) propogateToFollowers(data string) {
	fmt.Println("Propogating request to all followers..")
	followers := s.followerStore.GetAllKeys()

	for _, follower := range followers {
		_, err := follower.Write([]byte(fmt.Sprintf("%s\n", data)))
		if err != nil {
			fmt.Printf("Error writing to follower %s\n", err)
			continue
		}
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
