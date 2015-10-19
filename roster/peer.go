package roster

import "github.com/twstrike/coyim/xmpp"

// Peer represents and contains all the information you have about a specific peer.
// A Peer is always part of at least one roster.List, which is associated with an account.
type Peer struct {
	// The bare jid of the peer
	jid                string
	subscription       string
	name               string
	groups             map[string]bool
	status             string
	statusMsg          string
	offline            bool
	asked              bool
	pendingSubscribeId string
}

func toSet(ks ...string) map[string]bool {
	m := make(map[string]bool)
	for _, v := range ks {
		m[v] = true
	}
	return m
}

// PeerFrom returns a new Peer that contains the same information as the RosterEntry given
func PeerFrom(e xmpp.RosterEntry) *Peer {
	return &Peer{
		jid:          xmpp.RemoveResourceFromJid(e.Jid),
		subscription: e.Subscription,
		name:         e.Name,
		groups:       toSet(e.Group...),
	}
}

// PeerWithState returns a new Peer that contains the given state information
func PeerWithState(jid, status, statusMsg string) *Peer {
	return &Peer{
		jid:       xmpp.RemoveResourceFromJid(jid),
		status:    status,
		statusMsg: statusMsg,
	}
}

func merge(v1, v2 string) string {
	if v2 != "" {
		return v2
	}
	return v1
}

// MergeWith returns a new Peer that is the merger of the receiver and the argument, giving precedence to the argument when needed
func (p *Peer) MergeWith(p2 *Peer) *Peer {
	pNew := &Peer{}
	pNew.jid = p.jid
	pNew.subscription = merge(p.subscription, p2.subscription)
	pNew.name = merge(p.name, p2.name)
	pNew.status = merge(p.status, p2.status)
	pNew.statusMsg = merge(p.statusMsg, p2.statusMsg)
	pNew.offline = p2.offline
	pNew.asked = p2.asked
	pNew.pendingSubscribeId = merge(p.pendingSubscribeId, p2.pendingSubscribeId)
	pNew.groups = make(map[string]bool)
	for k, v := range p.groups {
		pNew.groups[k] = v
	}
	for k, v := range p2.groups {
		pNew.groups[k] = v
	}
	return pNew
}
