package roster

import (
	"fmt"
	"sort"
	"sync"

	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
)

// Status contains the status information for a specific resource
type Status struct {
	Status    string
	StatusMsg string
}

// Peer represents and contains all the information you have about a specific peer.
// A Peer is always part of at least one roster.List, which is associated with an account.
type Peer struct {
	Jid                jid.WithoutResource
	Subscription       string
	Name               string
	Nickname           string
	Groups             map[string]bool
	Asked              bool
	PendingSubscribeID string
	BelongsTo          string
	LatestError        *PeerError

	resources      map[string]Status
	resourcesLock  sync.RWMutex
	lockedResource jid.Resource

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
	return fmt.Sprintf("Peer{%s[%s (%s)], subscription='%s', status='%s'('%s') online=%v, asked=%v, pendingSubscribe='%s', belongsTo='%s', resources=%v, lockedResource='%s'}", p.Jid, p.Name, p.Nickname, p.Subscription, p.MainStatus(), p.MainStatusMsg(), p.IsOnline(), p.Asked, p.PendingSubscribeID, p.BelongsTo, p.resources, p.lockedResource)
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
		Jid:           jid.NR(e.Jid),
		Subscription:  e.Subscription,
		Name:          e.Name,
		Nickname:      nickname,
		Groups:        allGroups,
		HasConfigData: groups != nil && len(groups) > 0,
		BelongsTo:     belongsTo,
		Asked:         e.Ask == "subscribe",
		resources:     make(map[string]Status),
	}
}

// ToEntry returns a new RosterEntry with the same values
func (p *Peer) ToEntry() data.RosterEntry {
	return data.RosterEntry{
		Jid:          p.Jid.String(),
		Subscription: p.Subscription,
		Name:         p.Name,
		Group:        fromSet(p.Groups),
	}
}

// PeerWithState returns a new Peer that contains the given state information
func PeerWithState(jid jid.WithoutResource, status, statusMsg, belongsTo string, resource jid.Resource) *Peer {
	res := &Peer{
		Jid:       jid,
		BelongsTo: belongsTo,
		resources: make(map[string]Status),
	}
	res.AddResource(resource, status, statusMsg)
	return res
}

func peerWithPendingSubscribe(jid jid.WithoutResource, id, belongsTo string) *Peer {
	return &Peer{
		Jid:                jid,
		PendingSubscribeID: id,
		Asked:              true,
		BelongsTo:          belongsTo,
		resources:          make(map[string]Status),
	}
}

func merge(v1, v2 string) string {
	if v2 != "" {
		return v2
	}
	return v1
}

func union(v1, v2 map[string]Status) map[string]Status {
	if v1 == nil {
		v1 = make(map[string]Status)
	}
	if v2 == nil {
		v2 = make(map[string]Status)
	}
	res := make(map[string]Status)
	for r, v := range v1 {
		res[r] = v
	}
	for r, v := range v2 {
		res[r] = v
	}
	return res
}

// MergeWith returns a new Peer that is the merger of the receiver and the argument, giving precedence to the argument when needed
func (p *Peer) MergeWith(p2 *Peer) *Peer {
	pNew := &Peer{}
	pNew.Jid = p.Jid
	pNew.Subscription = merge(p.Subscription, p2.Subscription)
	pNew.Name = merge(p.Name, p2.Name)
	pNew.Nickname = merge(p.Nickname, p2.Nickname)
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
	return merge(p.Jid.String(), name)
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
func (p *Peer) AddResource(ss jid.Resource, status, statusMsg string) {
	s := string(ss)
	if s != "" {
		p.resourcesLock.Lock()
		defer p.resourcesLock.Unlock()

		p.resources[s] = Status{status, statusMsg}
	}
}

// RemoveResource removes the given resource
func (p *Peer) RemoveResource(s jid.Resource) {
	p.resourcesLock.Lock()
	defer p.resourcesLock.Unlock()

	delete(p.resources, string(s))
}

// Resources returns the resources for this peer
func (p *Peer) Resources() []jid.Resource {
	p.resourcesLock.RLock()
	defer p.resourcesLock.RUnlock()

	result1 := []string{}
	for k := range p.resources {
		result1 = append(result1, k)
	}
	sort.Strings(result1)

	result2 := []jid.Resource{}
	for _, k := range result1 {
		result2 = append(result2, jid.Resource(k))
	}

	return result2
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

	p.resources = make(map[string]Status)
}

// LastSeen sets the last resource used
func (p *Peer) LastSeen(r jid.Any) {
	p.lockedResource = r.PotentialResource()
}

// ResourceToUse returns the resource to use for this peer
func (p *Peer) ResourceToUse() jid.Resource {
	return p.lockedResource
}

// IsOnline returns true if any of the resources are online
func (p *Peer) IsOnline() bool {
	return len(p.resources) > 0
}

func (p *Peer) currentResourceStatus() Status {
	if p.lockedResource != jid.Resource("") {
		return p.resources[string(p.lockedResource)]
	}
	for _, s := range p.resources {
		return s
	}

	return Status{}
}

// MainStatus returns the status of the current main resource
func (p *Peer) MainStatus() string {
	return p.currentResourceStatus().Status
}

// MainStatusMsg returns the status message of the current main resource
func (p *Peer) MainStatusMsg() string {
	return p.currentResourceStatus().StatusMsg
}
