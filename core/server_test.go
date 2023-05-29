package core

import (
	"errors"
	"io"
	"testing"

	"gitub.com/sriramr98/go_kvdb/core/network"
	"gitub.com/sriramr98/go_kvdb/core/protocol"
	"gitub.com/sriramr98/go_kvdb/core/store"
)

type MockConn struct {
	io.ReadWriteCloser
	ReadFunc  func([]byte) (int, error)
	WriteFunc func([]byte) (int, error)
	CloseFunc func() error
}

func (mc *MockConn) Read(b []byte) (int, error) {
	if mc.ReadFunc != nil {
		return mc.ReadFunc(b)
	}
	return 0, errors.New("Read not implemented")
}

func (mc *MockConn) Write(b []byte) (int, error) {
	if mc.WriteFunc != nil {
		return mc.WriteFunc(b)
	}
	return 0, errors.New("Write not implemented")
}

func (mc *MockConn) Close() error {
	if mc.CloseFunc != nil {
		return mc.CloseFunc()
	}
	return errors.New("Close not implemented")
}

type MockDialer struct {
}

func (md *MockDialer) Dial(network, address string) (network.Conn, error) {
	return &MockConn{}, nil
}

type MockListener struct {
	AcceptFunc func() (network.Conn, error)
	CloseFunc  func() error
}

func (ml *MockListener) Accept() (network.Conn, error) {
	if ml.AcceptFunc != nil {
		return ml.AcceptFunc()
	}
	return nil, errors.New("Accept not implemented")
}

func (ml *MockListener) Close() error {
	if ml.CloseFunc != nil {
		return ml.CloseFunc()
	}
	return errors.New("Close not implemented")
}

func MockListenFunc(network, laddr string) (network.Listener, error) {
	return &MockListener{}, nil
}

type MockProcessor struct{}

func (mp *MockProcessor) Process(request protocol.Request) (protocol.Response, error) {
	if request.Command.Op != "SYNC" {
		return protocol.Response{}, errors.New("Expected SYNC command")
	}
	return protocol.Response{Success: true}, nil
}

type MockStore struct {
	store.DataStorer[network.Conn, struct{}]
}

type MockProtocol struct {
	protocol.Protocol
}

func TestNewServer(t *testing.T) {
	t.Run("Test New Server creation", func(t *testing.T) {
		opts := ServerOpts{
			IsLeader: true,
		}

		server, err := NewServer(
			opts,
			&MockProcessor{},
			&MockStore{},
			&MockProtocol{},
			&MockDialer{},
			MockListenFunc)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
			return
		}

		if server == nil {
			t.Errorf("Expected server to be non nil")
			return
		}

		if server.leaderConn != nil {
			t.Errorf("Expected leaderConn to be nil")
		}
	})

	t.Run("Test New Server for Follower", func(t *testing.T) {
		opts := ServerOpts{
			IsLeader: false,
		}

		server, err := NewServer(
			opts,
			&MockProcessor{},
			&MockStore{},
			&MockProtocol{},
			&MockDialer{},
			MockListenFunc)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
			return
		}

		if server == nil {
			t.Errorf("Expected server to be non nil")
			return
		}

		if server.leaderConn == nil {
			t.Errorf("Expected leaderConn to be non nil")
		}
	})

}
