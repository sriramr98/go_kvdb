package core

import (
	"errors"
	"strconv"
	"time"

	"gitub.com/sriramr98/go_kvdb/store"
)

var (
	ErrExpiryNotANumber   = errors.New("expiry value expected to be a number")
	ErrExpiryLessThanZero = errors.New("expiry seconds should be greater than 0")
)

type RequestProcessor interface {
	Process(request Request) (Response, error)
}

type CommandProcessor struct {
	Store *store.InMemoryStore
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

	if res.IsExpired(time.Now()) {
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
			return Response{}, ErrExpiryNotANumber
		}

		if ttlSet < 0 {
			return Response{}, ErrExpiryLessThanZero
		}
		ttlSeconds = ttlSet
	}

	expiresAt := -1
	if ttlSeconds > 0 {
		expiresAt = int(time.Now().Add(time.Duration(ttlSeconds) * time.Second).Unix())
	}

	cp.Store.Set(key, store.NewDataObj([]byte(value), expiresAt))

	return Response{Success: true}, nil
}

func (cp *CommandProcessor) processDel(request Request) (Response, error) {
	cp.Store.Delete(request.Params[0])
	return Response{Success: true}, nil
}

func (cp *CommandProcessor) processPing(_ Request) (Response, error) {
	return Response{Success: true, Value: []byte("PONG")}, nil
}
