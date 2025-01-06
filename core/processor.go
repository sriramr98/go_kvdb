package core

import (
	"strconv"
	"time"

	"gitub.com/sriramr98/go_kvdb/store"
)

type RequestProcessor interface {
	Process(request Request) (Response, error)
}

type CommandProcessor struct {
	Store store.DataStorer[string, store.DataObj]
}

func (cp *CommandProcessor) Process(request Request) (Response, error) {
	switch request.Command {
	case CMDGet:
		return cp.processGet(request)
	case CMDSet:
		return cp.processSet(request)
	case CMDDel:
		return cp.processDel(request)
	case CMDPing:
		return cp.processPing(request)
	default:
		return Response{}, ErrInvalidCommand
	}
}

func (cp *CommandProcessor) processGet(request Request) (Response, error) {
	res, err := cp.Store.Get(request.Params[0])
	if err != nil {
		return Response{}, err
	}

	if res.ExpiresAt() > 0 && time.Now().After(time.Unix(int64(res.ExpiresAt()), 0)) {
		// if key is already expired, delete it
		cp.Store.Delete(request.Params[0])
		return Response{}, store.ErrKeyNotFound
	}

	return Response{Success: true, Value: res.Value()}, nil
}

func (cp *CommandProcessor) processSet(request Request) (Response, error) {
	key := request.Params[0]
	value := request.Params[1]

	ttlSeconds := 0

	if len(request.Params) > 2 {
		ttlSet, err := strconv.Atoi(request.Params[2])
		if err != nil {
			return Response{}, err
		}
		ttlSeconds = ttlSet
	}

	cp.Store.Set(key, store.NewDataObj([]byte(value), ttlSeconds))

	return Response{Success: true}, nil
}

func (cp *CommandProcessor) processDel(request Request) (Response, error) {
	cp.Store.Delete(request.Params[0])
	return Response{Success: true}, nil
}

func (cp *CommandProcessor) processPing(_ Request) (Response, error) {
	return Response{Success: true, Value: []byte("PONG")}, nil
}

func (cp *CommandProcessor) processExpiry(key string, ttl time.Duration) {
	<-time.After(ttl)
	cp.Store.Delete(key)
}
