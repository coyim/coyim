package roster

import (
	"sort"
	"strings"
)

// Group represents a grouping of accounts and groups
type Group struct {
	GroupName     string
	fullGroupName []string
	peers         []*Peer
	groups        map[string]*Group
}

func (g *Group) findOrCreateInside(path, pathTilNow []string) *Group {
	if len(path) == 0 {
		return g
	}

	nowPath := append(pathTilNow, path[0])
	gg, ok := g.groups[path[0]]
	if ok {
		return gg.findOrCreateInside(path[1:], nowPath)
	}

	gg = &Group{
		GroupName:     path[0],
		fullGroupName: nowPath,
		peers:         make([]*Peer, 0),
		groups:        make(map[string]*Group, 0),
	}
	g.groups[path[0]] = gg

	return gg.findOrCreateInside(path[1:], nowPath)
}

// AddTo will add the peers in the receiver to the given group top level
func (l *List) AddTo(topLevel *Group, delim string) {
	for _, p := range l.peers {
		hadGroup := false
		for g, b := range p.Groups {
			if b {
				hadGroup = true
				gs := strings.Split(g, delim)
				gg := topLevel.findOrCreateInside(gs, []string{})
				gg.peers = append(gg.peers, p)
			}
		}
		if !hadGroup {
			topLevel.peers = append(topLevel.peers, p)
		}
	}
}

// TopLevelGroup returns a new top level group
func TopLevelGroup() *Group {
	return &Group{
		GroupName:     "",
		fullGroupName: []string{},
		peers:         make([]*Peer, 0),
		groups:        make(map[string]*Group, 0),
	}
}

// Grouped will group the peers in the given list
func (l *List) Grouped(delim string) *Group {
	topLevel := TopLevelGroup()

	l.AddTo(topLevel, delim)

	return topLevel
}

// Peers returns a sorted list of all the peers in this group
func (g *Group) Peers() []*Peer {
	sort.Sort(byNameForPresentation(g.peers))
	return g.peers
}

// UnsortedPeers returns an unsorted list of all the peers in this group
func (g *Group) UnsortedPeers() []*Peer {
	return g.peers
}

type byGroupNameAlphabetic []*Group

func (s byGroupNameAlphabetic) Len() int           { return len(s) }
func (s byGroupNameAlphabetic) Less(i, j int) bool { return s[i].GroupName < s[j].GroupName }
func (s byGroupNameAlphabetic) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// Groups returns a sorted list of all the groups in this group
func (g *Group) Groups() []*Group {
	allGroups := make([]*Group, 0, len(g.groups))
	for _, gg := range g.groups {
		allGroups = append(allGroups, gg)
	}

	sort.Sort(byGroupNameAlphabetic(allGroups))
	return allGroups
}

// FullGroupName returns the full group name
func (g *Group) FullGroupName() string {
	return strings.Join(g.fullGroupName, "::")
}
