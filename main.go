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
	flag.Parse()

	startDatabase(*isLeader)

	// initiateGracefulShutdown(ctx, clientServer, followerServer)

}

func startDatabase(isLeader bool) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	replicationStore := store.NewReplicationStore()
	inmemstore := store.NewInMemoryStore()
	processor := &processors.CommandProcessor{Store: inmemstore}

	runClientServer(replicationStore, processor, ctx, isLeader)
	if isLeader {
		// Only a leader can initiate replication to followers
		runReplicationServer(replicationStore, inmemstore, ctx)
	}

	select {}
}

func runClientServer(replicationStore store.DataStorer[net.Conn, struct{}], processor processors.RequestProcessor, ctx context.Context, isLeader bool) *core.Server {
	protocol := protocol.ClientProtocol{}
	opts := core.ServerOpts{Port: 8082, Role: core.ClientServerRole, IsLeader: isLeader, LeaderAddr: "localhost:8081"}

	clientServer := core.NewServer(opts, processor, replicationStore, protocol, ctx)
	go clientServer.Start()

	return clientServer
}

func runReplicationServer(replicationStore store.DataStorer[net.Conn, struct{}], clientStore store.DataStorer[string, []byte], ctx context.Context) *core.Server {
	protocol := protocol.FollowerProtocol{}
	opts := core.ServerOpts{Port: 8081, Role: core.ReplicaServerRole, IsLeader: true}

	processor := &processors.ReplicaProcessor{Store: clientStore}
	replicationServer := core.NewServer(opts, processor, replicationStore, protocol, ctx)
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
