package roster

import "github.com/twstrike/coyim/xmpp"

// List represent a list of peers. It takes care of both roster and presence information
// transparently and presents a unified view of this information to any UI
// List is not ordered, but can be asked to present its information in various orders
// depending on what policy is in use
// It also contains information about pending subscribes
// One invariant is that the list will only ever contain one Peer for each bare jid.
type List struct {
	peers map[string]*Peer
}

// New returns a new List
func New() *List {
	return &List{
		peers: make(map[string]*Peer),
	}
}

func (l *List) get(jid string) (*Peer, bool) {
	v, ok := l.peers[xmpp.RemoveResourceFromJid(jid)]
	return v, ok
}

// Remove returns the Peer with the jid from the List - it will first turn the jid into a bare jid.
// It returns true if it could remove the entry and false otherwise. It also returns the removed entry.
func (l *List) Remove(jid string) (*Peer, bool) {
	j := xmpp.RemoveResourceFromJid(jid)

	if v, ok := l.peers[j]; ok {
		delete(l.peers, j)
		return v, true
	}

	return nil, false
}

// AddOrMerge will add a new entry or merge with an existing entry the information from the given Peer
// It returns true if it added the entry and false otherwise
func (l *List) AddOrMerge(p *Peer) bool {
	if v, existed := l.peers[p.jid]; existed {
		l.peers[p.jid] = v.MergeWith(p)
		return false
	}

	l.peers[p.jid] = p

	return true
}

// AddOrReplace will add a new entry or replace an existing entry with the information from the given Peer
// It returns true if it added the entry and false otherwise
func (l *List) AddOrReplace(p *Peer) bool {
	_, existed := l.get(p.jid)

	l.peers[p.jid] = p

	return !existed
}

// PeerBecameUnavailable marks the peer as unavailable if they exist
// Returns true if they existed, otherwise false
func (l *List) PeerBecameUnavailable(jid string) bool {
	_, ok := l.Remove(jid)

	return ok
}

// PeerPresenceUpdate updates the status for the peer
// It returns true if it actually updated the status of the user
func (l *List) PeerPresenceUpdate(jid, status, statusMsg string) bool {
	if p, ok := l.get(jid); ok {
		if p.status != status || p.statusMsg != statusMsg {
			p.status = status
			p.statusMsg = statusMsg
			return true
		}
	} else {
		if status != "away" && status != "xa" {
			l.AddOrMerge(PeerWithState(jid, status, statusMsg))
			return true
		}
	}
	return false
}

// StateOf returns the status and status msg of the peer if it exists. It returns not ok if the peer doesn't exist.
func (l *List) StateOf(jid string) (status, statusMsg string, ok bool) {
	if p, existed := l.get(jid); existed {
		return p.status, p.statusMsg, true
	}

	return "", "", false
}

// - client presence
