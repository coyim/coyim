package muc

import "fmt"

// Affiliation represents an affiliation as specificed by section 5.2 in XEP-0045
type Affiliation interface {
	// IsBanned will return true if this specific affiliation means that the jid is banned from the room
	IsBanned() bool
	// IsMember returns true if this specific affiliation means that the jid is a member of the room
	IsMember() bool
	// IsModerator returns true if this specific affiliation means that the jid is a moderator of the room
	IsModerator() bool
}

type noneAffiliation struct{}
type outcastAffiliation struct{}
type memberAffiliation struct{}
type adminAffiliation struct{}
type ownerAffiliation struct{}

func (*noneAffiliation) IsBanned() bool    { return false }
func (*outcastAffiliation) IsBanned() bool { return true }
func (*memberAffiliation) IsBanned() bool  { return false }
func (*adminAffiliation) IsBanned() bool   { return false }
func (*ownerAffiliation) IsBanned() bool   { return false }

func (*noneAffiliation) IsMember() bool    { return false }
func (*outcastAffiliation) IsMember() bool { return false }
func (*memberAffiliation) IsMember() bool  { return true }
func (*adminAffiliation) IsMember() bool   { return true }
func (*ownerAffiliation) IsMember() bool   { return true }

func (*noneAffiliation) IsModerator() bool    { return false }
func (*outcastAffiliation) IsModerator() bool { return false }
func (*memberAffiliation) IsModerator() bool  { return false }
func (*adminAffiliation) IsModerator() bool   { return true }
func (*ownerAffiliation) IsModerator() bool   { return true }

// AffiliationFromString returns an Affiliation from the given string, or an error if the string doesn't match a known affiliation type
func AffiliationFromString(s string) (Affiliation, error) {
	switch s {
	case "none":
		return &noneAffiliation{}, nil
	case "outcast":
		return &outcastAffiliation{}, nil
	case "member":
		return &memberAffiliation{}, nil
	case "admin":
		return &adminAffiliation{}, nil
	case "owner":
		return &ownerAffiliation{}, nil
	default:
		return nil, fmt.Errorf("unknown affiliation string: '%s'", s)
	}
}
