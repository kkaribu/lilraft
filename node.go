package lilraft

const (
	stateFollower = iota
	stateCandidate
	stateLeader
)

// Log ...
type Log struct {
	id string

	// State
	state       int
	currentTerm uint64
	votedFor    string
	log         map[uint64]entry

	// Cluster log
	lastCommitted uint64
	lastApplied   uint64

	// Current term
	nextIndex  map[string]uint64
	matchIndex map[string]uint64
}

func (l *Log) appendEntries(term uint64, leaderID string, prevLogIndex, prevLogTerm uint64, entries []entry, leaderLastCommitted uint64) (uint64, bool) {
	// If it's from an older term, ignore it.
	if term < l.currentTerm {
		return 0, false
	}

	// The new entries can't be appended after a certain index if what the calling node
	// has at that index is different.
	if entry, ok := l.log[prevLogIndex]; !ok || entry.term != prevLogTerm {
		return 0, false
	}

	// If some of the new entries already exist in the node, their terms need to be checked.
	// As soon as both a local and a new entry have the same index but different terms, the
	// entry is discarded and so is the rest of the new entries.
	confirmedEntries := make([]entry, 0, len(entries))
	for i := range entries {
		if entry, ok := l.log[entries[i].index]; ok && entry.term != entries[i].term {
			break
		}
		confirmedEntries = append(confirmedEntries, entries[i])
	}

	// Append the new entries to the node's log.
	for i := range confirmedEntries {
		l.log[confirmedEntries[i].index] = confirmedEntries[i]
	}

	if len(confirmedEntries) == 0 {
		return l.currentTerm, true
	}

	// If the index of the last committed entry from the calling node (leader) is higher
	// then the index of the last committed entry locally, then the index of the last
	// committed entry becomes whatever is smaller between the leader's index and the index
	// of the last entry from those who just got appended.
	// In other words, the calling node considers itself the leader, otherwise it wouldn't
	// be attempting to append entries to this node. leaderLastCommitted is the leader's
	// highest index in its committed log. If this node's highest index is smaller, it
	// will have to be increased. It can't just be set to leaderLastCommitted, because maybe
	// some of the new entries from the leader haven't been committed yet.
	if l.lastCommitted < leaderLastCommitted {
		if leaderLastCommitted < confirmedEntries[len(confirmedEntries)-1].index {
			l.lastCommitted = leaderLastCommitted
		} else {
			l.lastCommitted = confirmedEntries[len(confirmedEntries)-1].index
		}
	}

	// Since the calling node considers itself the leader and has entries, the node must be
	// a follower.
	l.state = stateFollower

	return l.currentTerm, true
}

func (l *Log) requestVote(term uint64, candidateID string, lastLogIndex, lastLogTerm uint64) (uint64, bool) {
	// If it's from an older term, ignore it.
	if term < l.currentTerm {
		return 0, false
	}

	// The node will only responds positive to the vote request if it hasn't already voted for
	// a node or it has already voted for the calling node. Otherwise, it can't vote for more
	// than one node.
	// Also, a vote will only be granted if the calling node is at least up-to-date with this
	// node because if it is not, then for sure the calling node doesn't have what it takes to
	// be a leader.
	if l.votedFor == "" || l.votedFor == candidateID {
		if lastLogIndex >= l.lastCommitted {
			return l.currentTerm, true
		}
	}

	return 0, false
}

func (l *Log) isLeader() bool {
	return l.state == stateLeader
}
