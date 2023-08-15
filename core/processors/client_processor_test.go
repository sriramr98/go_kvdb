package processors

import (
	"testing"
	"time"

	"gitub.com/sriramr98/go_kvdb/core/protocol"
	"gitub.com/sriramr98/go_kvdb/mocks"
)

func TestCommandProcessor(t *testing.T) {

	store := mocks.NewDataStorer[string, []byte](t)
	cmd_processor := CommandProcessor{Store: store}

	t.Run("Test Get", func(t *testing.T) {
		request := protocol.Request{Command: protocol.CMDGet, Params: []string{"test"}}
		store.On("Get", "test").Return([]byte("test"), nil)
		cmd_processor.Process(request)

		store.AssertExpectations(t)
	})

	t.Run("Test Set", func(t *testing.T) {
		request := protocol.Request{Command: protocol.CMDSet, Params: []string{"test", "test"}}
		store.On("Set", "test", []byte("test")).Return(nil)
		cmd_processor.Process(request)

		store.AssertExpectations(t)
	})

	t.Run("Test Delete", func(t *testing.T) {
		request := protocol.Request{Command: protocol.CMDDel, Params: []string{"test"}}
		store.On("Delete", "test").Return(nil)
		cmd_processor.Process(request)

		store.AssertExpectations(t)
	})

	t.Run("Test Set with TTL", func(t *testing.T) {
		request := protocol.Request{Command: protocol.CMDSet, Params: []string{"test", "test", "2"}}
		store.On("Set", "test", []byte("test")).Return(nil)
		store.On("Delete", "test").Return(nil)

		cmd_processor.Process(request)
		time.Sleep(2 * time.Second)

		store.AssertExpectations(t)
	})
}
