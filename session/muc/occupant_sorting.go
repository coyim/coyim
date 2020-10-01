package muc

import (
	"fmt"
	"strings"
)

// ByOccupantNick sorts occupants by nickname
type ByOccupantNick []*Occupant

func (s ByOccupantNick) Len() int { return len(s) }
func (s ByOccupantNick) Less(i, j int) bool {
	return strings.ToLower(s[i].Nickname) < strings.ToLower(s[j].Nickname)
}
func (s ByOccupantNick) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// ByOccupantJid sorts occupants by real JID and then nickname
type ByOccupantJid []*Occupant

func (s ByOccupantJid) Len() int { return len(s) }
func (s ByOccupantJid) Less(i, j int) bool {
	differentJids := stringOrEmpty(s[i].RealJid) != stringOrEmpty(s[j].RealJid)
	if differentJids {
		return stringOrEmpty(s[i].RealJid) < stringOrEmpty(s[j].RealJid)
	}
	return s[i].Nickname < s[j].Nickname
}
func (s ByOccupantJid) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func stringOrEmpty(s fmt.Stringer) string {
	if s == nil {
		return ""
	}
	return s.String()
}
