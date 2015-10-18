package gui

import (
	"sort"

	"github.com/twstrike/coyim/xmpp"
)

type peer struct {
	jid          string
	subscription string
	name         string
	groups       []string
	status       string
	statusMsg    string
	offline      bool
}

type peers struct {
	peers []*peer
}

func fromRosterEntry(r xmpp.RosterEntry) *peer {
	p := &peer{}
	p.updateFromRosterEntry(r)
	return p
}

func (p *peer) updateFromRosterEntry(r xmpp.RosterEntry) {
	p.jid = r.Jid
	p.subscription = r.Subscription
	p.name = r.Name
	p.groups = r.Group
}

func (p *peer) shouldDisplay() bool {
	return p.subscription != "none" && p.subscription != ""
}

func (p *peer) subscribed() {
	switch p.subscription {
	case "from":
		p.subscription = "both"
	case "none", "":
		p.subscription = "to"
	}
}

func (p *peers) updateFromRosterEntries(r []xmpp.RosterEntry) {
	for _, re := range r {
		p.getOrAdd(re.Jid).updateFromRosterEntry(re)
	}
}

func fromRosterEntries(r []xmpp.RosterEntry) *peers {
	res := make([]*peer, len(r))
	for ix, e := range r {
		res[ix] = fromRosterEntry(e)
	}
	return createPeersFrom(res)
}

func createPeersFrom(ps []*peer) *peers {
	return &peers{ps}
}

func (p *peer) nameForPresentation() string {
	if p.name != "" {
		return p.name
	}
	return p.jid
}

func (p *peers) indexOf(jid string) (int, bool) {
	for ix, pp := range p.peers {
		if pp.jid == jid {
			return ix, true
		}
	}
	return 0, false
}

func (p *peers) getOrAdd(jid string) *peer {
	ix, ok := p.indexOf(jid)
	if ok {
		return p.peers[ix]
	}

	pr := &peer{jid: jid}
	p.peers = append(p.peers, pr)
	return pr
}

func (p *peers) remove(jid string) {
	ix, ok := p.indexOf(jid)
	if ok {
		p.peers = append(p.peers[:ix], p.peers[ix+1:]...)
	}
}

type byJid []*peer

func (s byJid) Len() int { return len(s) }
func (s byJid) Less(i, j int) bool {
	return s[i].jid < s[j].jid
}
func (s byJid) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (p *peers) iter(cb func(int, *peer)) {
	sort.Sort(byJid(p.peers))
	ix := 0
	for _, pr := range p.peers {
		if pr.shouldDisplay() {
			cb(ix, pr)
			ix++
		}
	}
}
