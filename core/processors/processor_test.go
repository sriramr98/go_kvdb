package processors

import (
	"testing"
	"time"

	"gitub.com/sriramr98/go_kvdb/core/protocol"
)

type MockDataStore struct {
	GetCalled    bool
	SetCalled    bool
	DeleteCalled bool
}

func (mds *MockDataStore) Get(key string) ([]byte, error) {
	mds.GetCalled = true
	return []byte("test"), nil
}

func (mds *MockDataStore) Set(key string, value []byte) {
	mds.SetCalled = true
}

func (mds *MockDataStore) Delete(key string) {
	mds.DeleteCalled = true
}

func TestCommandProcessor(t *testing.T) {

	store := &MockDataStore{}
	cmd_processor := CommandProcessor{Store: store}

	t.Run("Test Get", func(t *testing.T) {
		request := protocol.Request{Command: protocol.CMDGet, Params: []string{"test"}}
		cmd_processor.Process(request)
		if !store.GetCalled {
			t.Errorf("Get not called")
		}
	})

	t.Run("Test Set", func(t *testing.T) {
		request := protocol.Request{Command: protocol.CMDSet, Params: []string{"test", "test"}}
		cmd_processor.Process(request)
		if !store.SetCalled {
			t.Errorf("Set not called")
		}
	})

	t.Run("Test Delete", func(t *testing.T) {
		request := protocol.Request{Command: protocol.CMDDel, Params: []string{"test"}}
		cmd_processor.Process(request)
		if !store.DeleteCalled {
			t.Errorf("Delete not called")
		}
	})

	t.Run("Test Set with TTL", func(t *testing.T) {
		request := protocol.Request{Command: protocol.CMDSet, Params: []string{"test", "test", "2"}}
		cmd_processor.Process(request)
		if !store.SetCalled {
			t.Errorf("Set not called")
		}
		time.Sleep(2 * time.Second)
		if !store.DeleteCalled {
			t.Errorf("Delete not called for TTL")
		}
	})

}
