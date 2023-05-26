package processors

import (
	"gitub.com/sriramr98/go_kvdb/core/protocol"
)

type FollowerProcessor struct {
}

func (fp *FollowerProcessor) Process(request protocol.Request) (protocol.Response, error) {
	switch request.Command {
	case protocol.CMDSync:
		return fp.processSync(request)
	default:
		return protocol.Response{}, protocol.ErrInvalidCommand
	}
}

func (fp *FollowerProcessor) processSync(request protocol.Request) (protocol.Response, error) {
	return protocol.Response{Success: true}, nil
}
