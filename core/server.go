package core

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"

	"gitub.com/sriramr98/go_kvdb/core/processors"
	"gitub.com/sriramr98/go_kvdb/core/protocol"
	"gitub.com/sriramr98/go_kvdb/store"
)

type ServerRole int

const (
	ClientServerRole  ServerRole = 1
	ReplicaServerRole ServerRole = 2
	BufferSize                   = 2048
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

func NewServer(opts ServerOpts, clientProcessor processors.RequestProcessor, followerStore store.DataStorer[net.Conn, struct{}], protocol protocol.Protocol, ctx context.Context) (*Server, error) {
	server := &Server{opts: opts, requestProcessor: clientProcessor, protocol: protocol, followerStore: followerStore}

	if !opts.IsLeader {
		leaderConn, err := net.Dial("tcp", opts.LeaderAddr)
		if err != nil {
			return nil, fmt.Errorf("error connecting to leader: %w", err)
		}
		server.leaderConn = leaderConn
	}

	return server, nil
}

func (s *Server) Start() error {
	fmt.Println("Starting server on port", s.opts.Port)

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", s.opts.Port))
	if err != nil {
		return fmt.Errorf("error starting server: %w", err)
	}
	defer ln.Close()

	if !s.opts.IsLeader {
		if err := s.syncWithLeader(); err != nil {
			return err
		}
		go s.handleConnection(s.leaderConn, true)
	}

	return s.listenForConnections(ln)
}

func (s *Server) listenForConnections(ln net.Listener) error {
	for {
		conn, err := ln.Accept()
		if err != nil {
			return fmt.Errorf("error accepting connection: %w", err)
		}

		if s.opts.IsLeader && s.opts.Role == ReplicaServerRole {
			s.followerStore.Set(conn, struct{}{})
		}

		go s.handleConnection(conn, false)
	}
}

func (s *Server) syncWithLeader() error {
	// Send sync request to leader returns the current state of the leader
	_, err := s.leaderConn.Write([]byte("SYNC\n"))
	if err != nil {
		return fmt.Errorf("error syncing with leader: %w", err)
	}

	buf := make([]byte, BufferSize)
	n, err := s.leaderConn.Read(buf)
	if err != nil {
		return fmt.Errorf("error syncing with leader: %w", err)
	}

	data_received := buf[:n]
	parsedData := strings.SplitAfter(string(data_received), "OK: ")[1]
	fmt.Printf("Received state %s\n", parsedData)

	//TODO: Deserialize state
	//TODO: Set state in client store
	s.requestProcessor.Process(protocol.Request{
		Command: protocol.CMDSyncUpdate,
		Params:  []string{string(parsedData)},
	})

	return nil
}

func (s *Server) handleConnection(conn net.Conn, isFromLeader bool) {
	defer conn.Close()

	fmt.Println("New connection")

	for {
		buf := make([]byte, BufferSize)
		n, err := conn.Read(buf)
		if err != nil {
			s.handleReadError(err, conn)
			return
		}

		data_received := s.processReceivedData(buf, n)
		request, err := s.protocol.Parse(data_received)
		if err != nil {
			s.writeError(err, conn)
			continue
		}
		if !s.opts.IsLeader && !request.Command.CanFollowerProcess && !isFromLeader {
			s.writeError(fmt.Errorf("follower cannot process command %s", request.Command.Op), conn)
			return
		}
		if err := s.processCommand(request, isFromLeader, conn, data_received); err != nil {
			s.writeError(err, conn)
			continue
		}
	}
}

func (s *Server) processReceivedData(buf []byte, n int) string {
	data_received := string(buf[:n])
	data_received = data_received[:len(data_received)-1]
	fmt.Printf("Received %d bytes: %s\n", n, data_received)
	return data_received
}

func (s *Server) processCommand(request protocol.Request, isFromLeader bool, conn net.Conn, data_received string) error {
	response, err := s.requestProcessor.Process(request)
	if err != nil {
		return fmt.Errorf("error processing request: %w", err)
	}

	if s.opts.IsLeader && s.opts.Role == ClientServerRole && request.Command.IsReplicable {
		go s.propogateToFollowers(data_received)
	}

	if !isFromLeader {
		s.writeSuccess(response.Value, conn)
	}

	return nil
}

func (s *Server) handleReadError(err error, conn net.Conn) {
	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		fmt.Println("Connection timed out.. closing connection")
		return
	}

	if errors.Is(err, io.EOF) {
		if s.opts.IsLeader && s.opts.Role == ReplicaServerRole {
			fmt.Println("Follower closed connection")
			s.followerStore.Delete(conn)
		}
		fmt.Println("Connection closed")
		return
	}

	fmt.Printf("Error %s", err)
}

func (s *Server) propogateToFollowers(data string) {
	followers := s.followerStore.GetAllKeys()
	fmt.Printf("Propogating %s to %d followers", data, len(followers))

	for _, follower := range followers {
		_, err := follower.Write([]byte(fmt.Sprintf("%s\n", data)))
		if err != nil {
			fmt.Printf("Error writing to follower: %s\n", err)
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
