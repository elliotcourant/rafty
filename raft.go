package rafty

import (
	"fmt"
	"github.com/readystock/golog"
	"google.golang.org/grpc"
	"time"
)

type Raft struct {
	state raftState

	localGrpc *grpc.Server

	// RPC chan comes from the transport layer
	rpcChannel <-chan RPC
}

func (raft *Raft) setCurrentTerm(term uint64) {
	// persist to disk
	raft.state.setCurrentTerm(term)
}

func (raft *Raft) runFollower() {
	golog.Infof("entering follower state")
	heartbeatTimer := randomTimeout(time.Millisecond * 100)
	for {
		select {
		case rpc := <-raft.rpcChannel:
			raft.processRPC(rpc)
		}
	}
}

func (raft *Raft) processRPC(rpc RPC) {
	switch cmd := rpc.Command.(type) {
	case *AppendEntriesRequest:
		raft.appendEntries(rpc, cmd)
	case *RequestVoteRequest:

	default:
		golog.Errorf("got unexpected command: %#v", rpc.Command)
		rpc.Respond(nil, fmt.Errorf("unexpected command"))
	}
}

func (raft *Raft) appendEntries(rpc RPC, request *AppendEntriesRequest) {
	currentTerm := raft.state.currentTerm
	response := &AppendEntriesResponse{
		Term:           currentTerm,
		LastLog:        raft.state.lastLogIndex,
		Success:        false,
		NoRetryBackOff: false,
	}
	var rpcErr error
	defer func() {
		rpc.Respond(response, rpcErr)
	}()

	// Ignore an older term
	if request.Term < currentTerm {
		golog.Warnf("append entries request is from an older term")
		return
	}

	if request.Term > currentTerm || raft.state.state != Follower {
		raft.state.state = Follower
		raft.setCurrentTerm(request.Term)
		response.Term = request.Term
	}

}
