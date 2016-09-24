package roster

import (
	"sort"
	"sync"

	"github.com/twstrike/coyim/xmpp/utils"
)

// List represent a list of peers. It takes care of both roster and presence information
// transparently and presents a unified view of this information to any UI
// List is not ordered, but can be asked to present its information in various orders
// depending on what policy is in use
// It also contains information about pending subscribes
// One invariant is that the list will only ever contain one Peer for each bare jid.
type List struct {
	peers     map[string]*Peer
	peersLock sync.RWMutex
}

// New returns a new List
func New() *List {
	return &List{
		peers: make(map[string]*Peer),
	}
}

// Get returns the peer if it's known and false otherwise
func (l *List) Get(jid string) (*Peer, bool) {
	l.peersLock.RLock()
	defer l.peersLock.RUnlock()

	v, ok := l.peers[utils.RemoveResourceFromJid(jid)]
	return v, ok
}

// Clear removes all current entries in the list
func (l *List) Clear() {
	l.peersLock.Lock()
	defer l.peersLock.Unlock()

	l.peers = make(map[string]*Peer)
}

// Remove returns the Peer with the jid from the List - it will first turn the jid into a bare jid.
// It returns true if it could remove the entry and false otherwise. It also returns the removed entry.
func (l *List) Remove(jid string) (*Peer, bool) {
	j := utils.RemoveResourceFromJid(jid)

	l.peersLock.Lock()
	defer l.peersLock.Unlock()

	if v, ok := l.peers[j]; ok {
		delete(l.peers, j)
		return v, true
	}

	return nil, false
}

// AddOrMerge will add a new entry or merge with an existing entry the information from the given Peer
// It returns true if it added the entry and false otherwise
func (l *List) AddOrMerge(p *Peer) bool {
	l.peersLock.Lock()
	defer l.peersLock.Unlock()

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

	l.peersLock.Lock()
	defer l.peersLock.Unlock()
	l.peers[p.Jid] = p

	return !existed
}

// PeerBecameUnavailable marks the peer as unavailable if they exist
// Returns true if they existed, otherwise false
func (l *List) PeerBecameUnavailable(jid string) bool {
	if p, exist := l.Get(jid); exist {
		oldOnline := p.Online

		resource := utils.ResourceFromJid(jid)
		if resource != "" {
			p.RemoveResource(resource)
			p.Online = p.HasResources()
		} else {
			p.ClearResources()
			p.Online = false
		}

		return oldOnline != p.Online
	}

	return false
}

// PeerPresenceUpdate updates the status for the peer
// It returns true if it actually updated the status of the user
func (l *List) PeerPresenceUpdate(jid, status, statusMsg, belongsTo string) bool {
	resource := utils.ResourceFromJid(jid)

	if p, ok := l.Get(jid); ok {
		oldOnline := p.Online
		p.Online = true
		p.AddResource(resource)
		if p.Status != status || p.StatusMsg != statusMsg {
			p.Status = status
			p.StatusMsg = statusMsg
			return true
		}
		return !oldOnline
	}

	l.AddOrMerge(PeerWithState(jid, status, statusMsg, belongsTo, resource))
	return true
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

// LatestError will set the latest error on the jid in question
func (l *List) LatestError(jid string, code, tp, more string) {
	if p, existed := l.Get(jid); existed {
		p.SetLatestError(code, tp, more)
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

type byNameForPresentation []*Peer

func (s byNameForPresentation) Len() int { return len(s) }
func (s byNameForPresentation) Less(i, j int) bool {
	return s[i].NameForPresentation() < s[j].NameForPresentation()
}
func (s byNameForPresentation) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (l *List) intoSlice(res []*Peer) []*Peer {
	for _, v := range l.peers {
		res = append(res, v)
	}

	return res
}

// ToSlice returns a slice of all the peers in this roster list
func (l *List) ToSlice() []*Peer {
	l.peersLock.RLock()
	defer l.peersLock.RUnlock()

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
		func() {
			l.peersLock.RLock()
			defer l.peersLock.RUnlock()
			res = l.intoSlice(res)
		}()
	}

	sort.Sort(byJidAlphabetic(res))

	for ix, pr := range res {
		cb(ix, pr)
	}
}

// GetGroupNames return all group names for this peer list.
func (l *List) GetGroupNames() map[string]bool {
	l.peersLock.RLock()
	defer l.peersLock.RUnlock()

	names := map[string]bool{}

	//TODO: Should not group separator be part of a Peer List?
	for _, peer := range l.peers {
		for group := range peer.Groups {
			names[group] = true
		}
	}

	return names
}
