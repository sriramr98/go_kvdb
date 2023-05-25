package main

import (
	"gitub.com/sriramr98/go_kvdb/core"
	"gitub.com/sriramr98/go_kvdb/core/protocol"
	"gitub.com/sriramr98/go_kvdb/store"
)

func main() {
	inmemstore := store.NewInMemoryStore()
	processor := &core.CommandProcessor{Store: inmemstore}
	protcol := protocol.ClientProtocol{}
	core.NewServer(core.ServerOpts{Port: 8080}, processor, protcol).Start()
}
