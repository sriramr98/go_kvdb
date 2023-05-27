package processors

import (
	"gitub.com/sriramr98/go_kvdb/core/protocol"
	"gitub.com/sriramr98/go_kvdb/store"
)

type ReplicaProcessor struct {
	Store store.DataStorer[string, []byte]
}

func (fp *ReplicaProcessor) Process(request protocol.Request) (protocol.Response, error) {
	switch request.Command {
	case protocol.CMDSync:
		return fp.processSync(request)
	default:
		return protocol.Response{}, protocol.ErrInvalidCommand
	}
}

func (fp *ReplicaProcessor) processSync(request protocol.Request) (protocol.Response, error) {
	return protocol.Response{Success: true, Value: []byte("Got You")}, nil
}
