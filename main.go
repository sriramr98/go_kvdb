package main

import (
	"gitub.com/sriramr98/go_kvdb/core"
)

func main() {
	core.NewServer(core.ServerOpts{Port: 8080}).Start()
}
