package processors

import (
	"fmt"
	"strconv"
	"time"

	"github.com/vmihailenco/msgpack"
	"gitub.com/sriramr98/go_kvdb/core/protocol"
	"gitub.com/sriramr98/go_kvdb/core/store"
)

type CommandProcessor struct {
	Store store.DataStorer[string, []byte]
}

func (cp *CommandProcessor) Process(request protocol.Request) (protocol.Response, error) {
	fmt.Println("Client Processing request", request.Command)
	switch request.Command {
	case protocol.CMDGet:
		return cp.processGet(request)
	case protocol.CMDSet:
		return cp.processSet(request)
	case protocol.CMDDel:
		return cp.processDel(request)
	case protocol.CMDPing:
		return cp.processPing(request)
	case protocol.CMDSyncUpdate:
		return cp.processSyncUpdate(request)
	default:
		return protocol.Response{}, protocol.ErrInvalidCommand
	}
}

func (cp *CommandProcessor) processGet(request protocol.Request) (protocol.Response, error) {
	res, err := cp.Store.Get(request.Params[0])
	if err != nil {
		return protocol.Response{}, err
	}

	return protocol.Response{Success: true, Value: res}, nil
}

func (cp *CommandProcessor) processSet(request protocol.Request) (protocol.Response, error) {

	key := request.Params[0]
	value := request.Params[1]

	cp.Store.Set(key, []byte(value))

	if len(request.Params) > 2 {
		ttlSet, err := strconv.Atoi(request.Params[2])
		if err != nil {
			// If we're not able to set the TTL, we should delete the key
			cp.Store.Delete(key)
			return protocol.Response{}, err
		}
		go cp.processExpiry(key, time.Duration(ttlSet)*time.Second)
	}

	return protocol.Response{Success: true}, nil
}

func (cp *CommandProcessor) processDel(request protocol.Request) (protocol.Response, error) {
	cp.Store.Delete(request.Params[0])
	return protocol.Response{Success: true}, nil
}

func (cp *CommandProcessor) processPing(request protocol.Request) (protocol.Response, error) {
	return protocol.Response{Success: true, Value: []byte("PONG")}, nil
}

func (cp *CommandProcessor) processExpiry(key string, ttl time.Duration) {
	<-time.After(ttl)
	cp.Store.Delete(key)
}

func (cp *CommandProcessor) processSyncUpdate(request protocol.Request) (protocol.Response, error) {
	data := []byte(request.Params[0])

	var decoded map[string][]byte
	err := msgpack.Unmarshal(data, &decoded)
	if err != nil {
		fmt.Printf("Error decoding: %s", err)
		return protocol.Response{}, err
	}

	fmt.Printf("Decoded Map Len: %d", len(decoded))

	cp.Store.SetAll(decoded)

	return protocol.Response{Success: true}, nil
}
