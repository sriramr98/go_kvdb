package main

import (
	"gitub.com/sriramr98/go_kvdb/core"
	"gitub.com/sriramr98/go_kvdb/store"
)

func main() {
	inmemstore := store.NewInMemoryStore()
	processor := &core.CommandProcessor{Store: inmemstore}
	core.NewServer(core.ServerOpts{Port: 8080}, processor).Start()
}
