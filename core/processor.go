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
	Store store.DataStorer[string, []byte]
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

	return Response{Success: true, Value: res}, nil
}

func (cp *CommandProcessor) processSet(request Request) (Response, error) {

	key := request.Params[0]
	value := request.Params[1]

	cp.Store.Set(key, []byte(value))

	if len(request.Params) > 2 {
		ttlSet, err := strconv.Atoi(request.Params[2])
		if err != nil {
			// If we're not able to set the TTL, we should delete the key
			cp.Store.Delete(key)
			return Response{}, err
		}
		go cp.processExpiry(key, time.Duration(ttlSet)*time.Second)
	}

	return Response{Success: true}, nil
}

func (cp *CommandProcessor) processDel(request Request) (Response, error) {
	cp.Store.Delete(request.Params[0])
	return Response{Success: true}, nil
}

func (cp *CommandProcessor) processPing(request Request) (Response, error) {
	return Response{Success: true, Value: []byte("PONG")}, nil
}

func (cp *CommandProcessor) processExpiry(key string, ttl time.Duration) {
	<-time.After(ttl)
	cp.Store.Delete(key)
}
