package lilraft

// message ...
type message interface {
	term() uint64
	from() string
}

// msgAppendEntries ...
type msgAppendEntries struct {
	_term               uint64 // term
	leaderID            string // from
	prevLogIndex        uint64
	prevLogTerm         uint64
	entries             []entry
	leaderLastCommitted uint64
}

func (m msgAppendEntries) term() uint64 {
	return m._term
}

func (m msgAppendEntries) from() string {
	return m.leaderID
}

type msgRequestVote struct {
	_term        uint64 // term
	candidateID  string // from
	lastLogIndex uint64
	lastLogTerm  uint64
}

func (m msgRequestVote) term() uint64 {
	return m._term
}

func (m msgRequestVote) from() string {
	return m.candidateID
}
