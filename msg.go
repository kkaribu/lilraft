package lilraft

// MsgAppendEntries ...
type MsgAppendEntries struct {
	term                uint64
	leaderID            string
	prevLogIndex        uint64
	prevLogTerm         uint64
	entries             []entry
	leaderLastCommitted uint64
}

// MsgRequestVote ...
type MsgRequestVote struct {
	term         uint64
	candidateID  string
	lastLogIndex uint64
	lastLogTerm  uint64
}

// Msg ...
type Msg interface{}
