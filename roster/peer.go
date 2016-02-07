package roster

import (
	"fmt"

	"github.com/twstrike/coyim/xmpp/data"
	xutils "github.com/twstrike/coyim/xmpp/utils"
)

// Peer represents and contains all the information you have about a specific peer.
// A Peer is always part of at least one roster.List, which is associated with an account.
type Peer struct {
	// The bare jid of the peer
	Jid                string
	Subscription       string
	Name               string
	Nickname           string
	Groups             map[string]bool
	Status             string
	StatusMsg          string
	Online             bool
	Asked              bool
	PendingSubscribeID string
	BelongsTo          string
	LatestError        *PeerError
}

// PeerError contains information about an error for this peer
type PeerError struct {
	Code string
	Type string
	More string
}

func toSet(ks ...string) map[string]bool {
	m := make(map[string]bool)
	for _, v := range ks {
		m[v] = true
	}
	return m
}

func fromSet(vs map[string]bool) []string {
	m := make([]string, 0, len(vs))
	for k, v := range vs {
		if v {
			m = append(m, k)
		}
	}
	return m
}

// Dump will dump all info in the peer in a very verbose format
func (p *Peer) Dump() string {
	return fmt.Sprintf("Peer{%s[%s (%s)], subscription='%s', status='%s'('%s') online=%v, asked=%v, pendingSubscribe='%s', belongsTo='%s'}", p.Jid, p.Name, p.Nickname, p.Subscription, p.Status, p.StatusMsg, p.Online, p.Asked, p.PendingSubscribeID, p.BelongsTo)
}

// PeerFrom returns a new Peer that contains the same information as the RosterEntry given
func PeerFrom(e data.RosterEntry, belongsTo, nickname string) *Peer {
	return &Peer{
		Jid:          xutils.RemoveResourceFromJid(e.Jid),
		Subscription: e.Subscription,
		Name:         e.Name,
		Nickname:     nickname,
		Groups:       toSet(e.Group...),
		BelongsTo:    belongsTo,
	}
}

// ToEntry returns a new RosterEntry with the same values
func (p *Peer) ToEntry() data.RosterEntry {
	return data.RosterEntry{
		Jid:          p.Jid,
		Subscription: p.Subscription,
		Name:         p.Name,
		Group:        fromSet(p.Groups),
	}
}

// PeerWithState returns a new Peer that contains the given state information
func PeerWithState(jid, status, statusMsg, belongsTo string) *Peer {
	return &Peer{
		Jid:       xutils.RemoveResourceFromJid(jid),
		Status:    status,
		StatusMsg: statusMsg,
		Online:    true,
		BelongsTo: belongsTo,
	}
}

func peerWithPendingSubscribe(jid, id, belongsTo string) *Peer {
	return &Peer{
		Jid:                xutils.RemoveResourceFromJid(jid),
		PendingSubscribeID: id,
		BelongsTo:          belongsTo,
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
	pNew.Jid = p.Jid
	pNew.Subscription = merge(p.Subscription, p2.Subscription)
	pNew.Name = merge(p.Name, p2.Name)
	pNew.Nickname = merge(p.Nickname, p2.Nickname)
	pNew.Status = merge(p.Status, p2.Status)
	pNew.StatusMsg = merge(p.StatusMsg, p2.StatusMsg)
	pNew.Online = p.Online || p2.Online
	pNew.Asked = p2.Asked
	pNew.PendingSubscribeID = merge(p.PendingSubscribeID, p2.PendingSubscribeID)
	pNew.Groups = make(map[string]bool)
	pNew.BelongsTo = merge(p.BelongsTo, p2.BelongsTo)
	if len(p2.Groups) > 0 {
		pNew.Groups = p2.Groups
	} else {
		pNew.Groups = p.Groups
	}
	return pNew
}

// NameForPresentation returns the name if it exists and otherwise the JID
func (p *Peer) NameForPresentation() string {
	name := merge(p.Name, p.Nickname)
	return merge(p.Jid, name)
}

// SetLatestError will set the latest error on the jid in question
func (p *Peer) SetLatestError(code, tp, more string) {
	p.LatestError = &PeerError{code, tp, more}
}
