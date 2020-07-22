package muc

import "fmt"

// ByOccupantNick sorts occupants by nickname
type ByOccupantNick []*Occupant

func (s ByOccupantNick) Len() int           { return len(s) }
func (s ByOccupantNick) Less(i, j int) bool { return s[i].Nick < s[j].Nick }
func (s ByOccupantNick) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// ByOccupantJid sorts occupants by real JID and then nickname
type ByOccupantJid []*Occupant

func (s ByOccupantJid) Len() int { return len(s) }
func (s ByOccupantJid) Less(i, j int) bool {
	differentJids := stringOrEmpty(s[i].Jid) != stringOrEmpty(s[j].Jid)
	if differentJids {
		return stringOrEmpty(s[i].Jid) < stringOrEmpty(s[j].Jid)
	}
	return s[i].Nick < s[j].Nick
}
func (s ByOccupantJid) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func stringOrEmpty(s fmt.Stringer) string {
	if s == nil {
		return ""
	}
	return s.String()
}
