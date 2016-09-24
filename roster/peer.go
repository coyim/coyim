package roster

import (
	"fmt"
	"sort"
	"sync"

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
	// Instead of bool here, we should be able to use priority for the resource, when we want
	resources     map[string]bool
	resourcesLock sync.RWMutex
	lastResource  string
	HasConfigData bool
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
	return fmt.Sprintf("Peer{%s[%s (%s)], subscription='%s', status='%s'('%s') online=%v, asked=%v, pendingSubscribe='%s', belongsTo='%s', resources=%v, lastResource='%s'}", p.Jid, p.Name, p.Nickname, p.Subscription, p.Status, p.StatusMsg, p.Online, p.Asked, p.PendingSubscribeID, p.BelongsTo, p.resources, p.lastResource)
}

// PeerFrom returns a new Peer that contains the same information as the RosterEntry given
func PeerFrom(e data.RosterEntry, belongsTo, nickname string, groups []string) *Peer {
	// merge remote and local groups
	g := groups
	if g == nil || len(g) == 0 {
		g = e.Group
	}
	allGroups := toSet(g...)

	return &Peer{
		Jid:           xutils.RemoveResourceFromJid(e.Jid),
		Subscription:  e.Subscription,
		Name:          e.Name,
		Nickname:      nickname,
		Groups:        allGroups,
		HasConfigData: groups != nil && len(groups) > 0,
		BelongsTo:     belongsTo,
		resources:     toSet(),
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
func PeerWithState(jid, status, statusMsg, belongsTo, resource string) *Peer {
	res := &Peer{
		Jid:       xutils.RemoveResourceFromJid(jid),
		Status:    status,
		StatusMsg: statusMsg,
		Online:    true,
		BelongsTo: belongsTo,
		resources: toSet(),
	}
	res.AddResource(resource)
	return res
}

func peerWithPendingSubscribe(jid, id, belongsTo string) *Peer {
	return &Peer{
		Jid:                xutils.RemoveResourceFromJid(jid),
		PendingSubscribeID: id,
		BelongsTo:          belongsTo,
		resources:          toSet(),
	}
}

func merge(v1, v2 string) string {
	if v2 != "" {
		return v2
	}
	return v1
}

func union(v1, v2 map[string]bool) map[string]bool {
	if v1 == nil {
		v1 = toSet()
	}
	if v2 == nil {
		v2 = toSet()
	}
	v1v := fromSet(v1)
	v2v := fromSet(v2)
	return toSet(append(v1v, v2v...)...)
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
	if p.HasConfigData || len(p2.Groups) == 0 {
		pNew.Groups = p.Groups
		pNew.HasConfigData = p.HasConfigData
	} else {
		pNew.Groups = p2.Groups
		pNew.HasConfigData = p2.HasConfigData
	}

	pNew.resources = union(p.resources, p2.resources)

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

// SetGroups set the Peer groups
func (p *Peer) SetGroups(groups []string) {
	p.Groups = toSet(groups...)
}

// AddResource adds the given resource if it isn't blank
func (p *Peer) AddResource(s string) {
	if s != "" {
		p.resourcesLock.Lock()
		defer p.resourcesLock.Unlock()

		p.resources[s] = true
	}
}

// RemoveResource removes the given resource
func (p *Peer) RemoveResource(s string) {
	p.resourcesLock.Lock()
	defer p.resourcesLock.Unlock()

	delete(p.resources, s)
}

// Resources returns the resources for this peer
func (p *Peer) Resources() []string {
	p.resourcesLock.RLock()
	defer p.resourcesLock.RUnlock()

	result := []string{}
	for k := range p.resources {
		result = append(result, k)
	}
	sort.Strings(result)

	return result
}

// HasResources returns true if this peer has any online resources
func (p *Peer) HasResources() bool {
	p.resourcesLock.RLock()
	defer p.resourcesLock.RUnlock()

	return len(p.resources) > 0
}

// ClearResources removes all known resources for the given peer
func (p *Peer) ClearResources() {
	p.resourcesLock.Lock()
	defer p.resourcesLock.Unlock()

	p.resources = toSet()
}

// LastResource sets the last resource used, if not empty
func (p *Peer) LastResource(r string) {
	if r != "" {
		p.lastResource = r
	}
}

// ResourceToUse returns the resource to use for this peer
func (p *Peer) ResourceToUse() string {
	if p.lastResource == "" {
		return ""
	}
	p.resourcesLock.RLock()
	defer p.resourcesLock.RUnlock()

	if p.resources[p.lastResource] {
		return p.lastResource
	}

	return ""
}
