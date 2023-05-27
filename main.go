package main

import (
	"context"
	"flag"
	"net"

	"gitub.com/sriramr98/go_kvdb/core"
	"gitub.com/sriramr98/go_kvdb/core/processors"
	"gitub.com/sriramr98/go_kvdb/core/protocol"
	"gitub.com/sriramr98/go_kvdb/store"
)

func main() {

	isLeader := flag.Bool("leader", false, "Start the server as a leader")
	mainPort := flag.Int("port", 8080, "Port to start the server on")
	flag.Parse()

	startDatabase(*isLeader, *mainPort)

	// initiateGracefulShutdown(ctx, clientServer, followerServer)

}

func startDatabase(isLeader bool, port int) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	replicationStore := store.NewReplicationStore()
	inmemstore := store.NewInMemoryStore()
	processor := &processors.CommandProcessor{Store: inmemstore}

	runClientServer(port, replicationStore, processor, ctx, isLeader)
	if isLeader {
		// Only a leader can initiate replication to followers
		runReplicationServer(replicationStore, inmemstore, ctx)
	}

	select {}
}

func runClientServer(port int, replicationStore store.DataStorer[net.Conn, struct{}], processor processors.RequestProcessor, ctx context.Context, isLeader bool) *core.Server {
	protocol := protocol.ClientProtocol{}
	opts := core.ServerOpts{Port: port, Role: core.ClientServerRole, IsLeader: isLeader, LeaderAddr: "localhost:8081"}

	clientServer, err := core.NewServer(opts, processor, replicationStore, protocol, ctx)
	if err != nil {
		panic(err)
	}
	go clientServer.Start()

	return clientServer
}

func runReplicationServer(replicationStore store.DataStorer[net.Conn, struct{}], clientStore store.DataStorer[string, []byte], ctx context.Context) *core.Server {
	protocol := protocol.FollowerProtocol{}
	opts := core.ServerOpts{Port: 8081, Role: core.ReplicaServerRole, IsLeader: true}

	processor := &processors.ReplicaProcessor{Store: clientStore}
	replicationServer, err := core.NewServer(opts, processor, replicationStore, protocol, ctx)
	if err != nil {
		panic(err)
	}
	go replicationServer.Start()

	return replicationServer
}

// func initiateGracefulShutdown(ctx context.Context, clientServer, replicationServer *core.Server) {

// 	stopCh := make(chan os.Signal, 1)
// 	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)

// 	<-stopCh

// 	var wg sync.WaitGroup
// 	wg.Add(2)

// 	go func() {
// 		defer wg.Done()
// 		clientServer.Stop(ctx)
// 	}()

// 	go func() {
// 		defer wg.Done()
// 		replicationServer.Stop(ctx)
// 	}()

// 	wg.Wait()
// }
