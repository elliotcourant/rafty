package rafty

import (
	"sync"
)

type RaftState uint32

const (
	// Follower is the state of an established node receiving it's state from the leader.
	Follower RaftState = iota

	// Candidate is the state of a node in transition.
	Candidate

	// Leader is the state of the node that dictates the state of the cluster.
	Leader

	// Shutdown is the terminal state of a raft node.
	Shutdown
)

type raftState struct {
	currentTerm uint64

	commitIndex uint64

	lastApplied uint64

	lastLock sync.Mutex

	lastSnapshotIndex uint64
	lastSnapshotTerm  uint64

	lastLogIndex uint64
	lastLogTerm  uint64

	state RaftState
}

func (state raftState) setCurrentTerm(term uint64) {
	state.currentTerm = term
}
