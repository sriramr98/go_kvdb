package processors

import (
	"fmt"

	"github.com/vmihailenco/msgpack"
	"gitub.com/sriramr98/go_kvdb/core/protocol"
	"gitub.com/sriramr98/go_kvdb/store"
)

type ReplicaProcessor struct {
	Store store.DataStorer[string, []byte]
}

func (fp *ReplicaProcessor) Process(request protocol.Request) (protocol.Response, error) {
	fmt.Println("Replica Processing request", request.Command)
	switch request.Command {
	case protocol.CMDSync:
		return fp.processSync(request)
	default:
		return protocol.Response{}, protocol.ErrInvalidCommand
	}
}

func (fp *ReplicaProcessor) processSync(request protocol.Request) (protocol.Response, error) {
	data := fp.Store.GetAll()

	fmt.Printf("Sending %d data through SYNC", len(data))

	encoded, err := msgpack.Marshal(data)
	fmt.Printf("Encoded %s", encoded)
	if err != nil {
		return protocol.Response{}, err
	}

	return protocol.Response{Success: true, Value: encoded}, nil
}
