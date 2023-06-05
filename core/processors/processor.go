package processors

import (
	"gitub.com/sriramr98/go_kvdb/core/protocol"
)

type Processor interface {
	Process(request protocol.Request) (protocol.Response, error)
}
