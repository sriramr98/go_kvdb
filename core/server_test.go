package core

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitub.com/sriramr98/go_kvdb/core/network"
	"gitub.com/sriramr98/go_kvdb/mocks"
)

func TestNewServer(t *testing.T) {

	t.Run("Test New Server creation for Leader", func(t *testing.T) {
		mockStore := mocks.NewDataStorer[network.Conn, struct{}](t)
		mockProtocol := mocks.NewProtocol(t)
		mockDialer := mocks.NewDialer(t)
		mockListenFunc := mocks.NewListener(t)
		mockProcessor := mocks.NewRequestProcessor(t)

		opts := ServerOpts{
			IsLeader: true,
			Port:     1234,
		}

		server, err := NewServer(
			opts,
			mockProcessor,
			mockStore,
			mockProtocol,
			mockDialer,
			mockListenFunc,
		)

		assert.NoError(t, err)
		assert.NotNil(t, server)
		assert.Nil(t, server.leaderConn)

		mockDialer.AssertExpectations(t)
		mockProcessor.AssertExpectations(t)
		mockStore.AssertExpectations(t)
		mockProtocol.AssertExpectations(t)
		mockListenFunc.AssertExpectations(t)

	})

	t.Run("Test New Server for Follower", func(t *testing.T) {
		opts := ServerOpts{
			IsLeader: false,
		}

		mockStore := mocks.NewDataStorer[network.Conn, struct{}](t)
		mockProtocol := mocks.NewProtocol(t)
		mockDialer := mocks.NewDialer(t)
		mockListenFunc := mocks.NewListener(t)
		mockProcessor := mocks.NewRequestProcessor(t)

		mockDialer.On("Dial", "tcp", mock.Anything).Return(mocks.NewConn(t), nil)

		server, err := NewServer(
			opts,
			mockProcessor,
			mockStore,
			mockProtocol,
			mockDialer,
			mockListenFunc,
		)

		assert.NoError(t, err)
		assert.NotNil(t, server)
		assert.NotNil(t, server.leaderConn)

		mockDialer.AssertExpectations(t)
		mockProcessor.AssertExpectations(t)
		mockStore.AssertExpectations(t)
		mockProtocol.AssertExpectations(t)
		mockListenFunc.AssertExpectations(t)
	})
}

func TestServer_Start(t *testing.T) {
	//TODO: Figure out how to test start since this runs an infinite for loop and never finishes for our assertions to work
	// t.Run("Successful Start", func(t *testing.T) {
	// 	mockListener.On("Listen", "tcp", mock.Anything).Return(mockNetProcessor, nil)
	// 	mockNetProcessor.On("Accept").Return(mockConn, nil)
	// 	mockNetProcessor.On("Close").Return(nil)

	// 	mockProcessor.On("Process", mock.Anything).Return(protocol.Response{}, nil)

	// 	mockConn.On("Write", mock.Anything).Return(0, nil)
	// 	mockConn.On("Read", mock.Anything).Return(0, errors.New("connection error"))
	// 	mockConn.On("Close").Return(nil)
	// 	err := server.Start()
	// 	assert.Error(t, err, "error accepting connection: connection error")
	// 	mockListener.AssertExpectations(t)
	// 	mockProcessor.AssertExpectations(t)
	// 	mockConn.AssertExpectations(t)
	// })

	t.Run("Error Starting Server", func(t *testing.T) {
		mockStore := mocks.NewDataStorer[network.Conn, struct{}](t)
		mockProtocol := mocks.NewProtocol(t)
		mockDialer := mocks.NewDialer(t)
		mockListener := mocks.NewListener(t)
		mockProcessor := mocks.NewRequestProcessor(t)
		mockConn := mocks.NewConn(t)

		// Create the server instance
		server := &Server{
			opts:             ServerOpts{Port: 1234, Role: ClientServerRole, IsLeader: true},
			requestProcessor: mockProcessor,
			protocol:         mockProtocol, // Replace with the actual protocol implementation
			followerStore:    mockStore,    // Replace with the actual follower store implementation
			dialer:           mockDialer,   // Replace with the actual dialer implementation
			listener:         mockListener,
		}

		// Set up expectations for the mock listener to return an error when Listen is called
		mockListener.On("Listen", "tcp", ":1234").Return(nil, errors.New("start error"))

		// Call the Start function
		err := server.Start()

		// Assertions
		assert.Error(t, err, "error starting server: start error")
		mockListener.AssertExpectations(t)
		mockProcessor.AssertExpectations(t)
		mockConn.AssertExpectations(t)
	})

	t.Run("Error Accepting Connection", func(t *testing.T) {
		mockStore := mocks.NewDataStorer[network.Conn, struct{}](t)
		mockProtocol := mocks.NewProtocol(t)
		mockDialer := mocks.NewDialer(t)
		mockListener := mocks.NewListener(t)
		mockProcessor := mocks.NewRequestProcessor(t)
		mockConn := mocks.NewConn(t)
		mockNetProcessor := mocks.NewProcessor(t)

		// Create the server instance
		server := &Server{
			opts:             ServerOpts{Port: 1234, Role: ClientServerRole, IsLeader: true},
			requestProcessor: mockProcessor,
			protocol:         mockProtocol, // Replace with the actual protocol implementation
			followerStore:    mockStore,    // Replace with the actual follower store implementation
			dialer:           mockDialer,   // Replace with the actual dialer implementation
			listener:         mockListener,
		}

		// Set up expectations for the mock listener to return the mock listener instance and an error when Accept is called
		mockListener.On("Listen", "tcp", ":1234").Return(mockNetProcessor, nil)
		mockNetProcessor.On("Accept").Return(nil, errors.New("accept error"))
		mockNetProcessor.On("Close").Return(nil)

		// Call the Start function
		err := server.Start()

		// Assertions
		assert.Error(t, err, "error accepting connection: accept error")
		mockListener.AssertExpectations(t)
		mockProcessor.AssertExpectations(t)
		mockConn.AssertExpectations(t)
	})

}
