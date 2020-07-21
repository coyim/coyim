package muc

// ByOccupantNick sorts occupants by nickname
type ByOccupantNick []*Occupant

func (s ByOccupantNick) Len() int           { return len(s) }
func (s ByOccupantNick) Less(i, j int) bool { return s[i].Nick < s[j].Nick }
func (s ByOccupantNick) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// ByOccupantJid sorts occupants by real JID and then nickname
type ByOccupantJid []*Occupant

func (s ByOccupantJid) Len() int { return len(s) }
func (s ByOccupantJid) Less(i, j int) bool {
	differentJids := s[i].Jid.String() != s[j].Jid.String()
	if differentJids {
		return s[i].Jid.String() < s[j].Jid.String()
	}
	return s[i].Nick < s[j].Nick
}
func (s ByOccupantJid) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
