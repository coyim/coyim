package roster

import (
	"sort"

	"github.com/twstrike/coyim/xmpp"
)

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

// Get returns the peer if it's known and false otherwise
func (l *List) Get(jid string) (*Peer, bool) {
	v, ok := l.peers[xmpp.RemoveResourceFromJid(jid)]
	return v, ok
}

// Clear removes all current entries in the list
func (l *List) Clear() {
	l.peers = make(map[string]*Peer)
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
	if v, existed := l.peers[p.Jid]; existed {
		l.peers[p.Jid] = v.MergeWith(p)
		return false
	}

	l.peers[p.Jid] = p

	return true
}

// AddOrReplace will add a new entry or replace an existing entry with the information from the given Peer
// It returns true if it added the entry and false otherwise
func (l *List) AddOrReplace(p *Peer) bool {
	_, existed := l.Get(p.Jid)

	l.peers[p.Jid] = p

	return !existed
}

// PeerBecameUnavailable marks the peer as unavailable if they exist
// Returns true if they existed, otherwise false
func (l *List) PeerBecameUnavailable(jid string) bool {
	if p, exist := l.Get(jid); exist {
		p.Online = false
		return true
	}

	return false
}

// PeerPresenceUpdate updates the status for the peer
// It returns true if it actually updated the status of the user
func (l *List) PeerPresenceUpdate(jid, status, statusMsg, belongsTo string) bool {
	if p, ok := l.Get(jid); ok {
		oldOnline := p.Online
		p.Online = true
		if p.Status != status || p.StatusMsg != statusMsg {
			p.Status = status
			p.StatusMsg = statusMsg
			return true
		}
		return !oldOnline
	}

	if status != "away" && status != "xa" {
		l.AddOrMerge(PeerWithState(jid, status, statusMsg, belongsTo))
		return true
	}

	return false
}

// StateOf returns the status and status msg of the peer if it exists. It returns not ok if the peer doesn't exist.
func (l *List) StateOf(jid string) (status, statusMsg string, ok bool) {
	if p, existed := l.Get(jid); existed {
		return p.Status, p.StatusMsg, true
	}

	return "", "", false
}

// SubscribeRequest adds a new pending subscribe request
func (l *List) SubscribeRequest(jid, id, belongsTo string) {
	l.AddOrMerge(peerWithPendingSubscribe(jid, id, belongsTo))
}

// RemovePendingSubscribe will return a subscribe id and remove the pending subscribe if it exists
// It will return false if no such subscribe is in flight
func (l *List) RemovePendingSubscribe(jid string) (string, bool) {
	if p, existed := l.Get(jid); existed {
		s := p.PendingSubscribeID
		p.PendingSubscribeID = ""
		return s, s != ""
	}

	return "", false
}

// GetPendingSubscribe will return a subscribe id without removing it
func (l *List) GetPendingSubscribe(jid string) (string, bool) {
	if p, existed := l.Get(jid); existed {
		return p.PendingSubscribeID, p.PendingSubscribeID != ""
	}

	return "", false
}

// Subscribed will mark the jid as subscribed
func (l *List) Subscribed(jid string) {
	if p, existed := l.Get(jid); existed {
		switch p.Subscription {
		case "from":
			p.Subscription = "both"
		case "none", "":
			p.Subscription = "to"
		}
		p.PendingSubscribeID = ""
		p.Asked = false
	}
}

// Unsubscribed will mark the jid as unsubscribed
func (l *List) Unsubscribed(jid string) {
	if p, existed := l.Get(jid); existed {
		switch p.Subscription {
		case "both":
			p.Subscription = "from"
		case "to":
			p.Subscription = "none"
		}
		p.PendingSubscribeID = ""
		p.Asked = false
	}
}

type byJidAlphabetic []*Peer

func (s byJidAlphabetic) Len() int           { return len(s) }
func (s byJidAlphabetic) Less(i, j int) bool { return s[i].Jid < s[j].Jid }
func (s byJidAlphabetic) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func (l *List) intoSlice(res []*Peer) []*Peer {
	for _, v := range l.peers {
		res = append(res, v)
	}

	return res
}

// ToSlice returns a slice of all the peers in this roster list
func (l *List) ToSlice() []*Peer {
	res := l.intoSlice(make([]*Peer, 0, len(l.peers)))

	sort.Sort(byJidAlphabetic(res))

	return res
}

// Iter calls the cb function once for each peer in the list
func (l *List) Iter(cb func(int, *Peer)) {
	for ix, pr := range l.ToSlice() {
		cb(ix, pr)
	}
}

// IterAll calls the cb function once for each peer in all the lists
func IterAll(cb func(int, *Peer), ls ...*List) {
	res := make([]*Peer, 0, 20)
	for _, l := range ls {
		res = l.intoSlice(res)
	}

	sort.Sort(byJidAlphabetic(res))

	for ix, pr := range res {
		cb(ix, pr)
	}
}
