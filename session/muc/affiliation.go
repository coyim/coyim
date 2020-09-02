package muc

import "fmt"

const (
	// AffiliationOwner represents XMPP muc 'owner' affiliation
	AffiliationOwner = "owner"
	// AffiliationAdmin represents XMPP muc 'admin' affiliation
	AffiliationAdmin = "admin"
	// AffiliationMember represents XMPP muc 'member' affiliation
	AffiliationMember = "member"
	// AffiliationOutcast represents XMPP muc 'outcast' affiliation
	AffiliationOutcast = "outcast"
	// AffiliationNone represents XMPP muc 'none' affiliation
	AffiliationNone = "none"
)

// Affiliation represents an affiliation as specificed by section 5.2 in XEP-0045
type Affiliation interface {
	// IsBanned will return true if this specific affiliation means that the jid is banned from the room
	IsBanned() bool
	// IsMember returns true if this specific affiliation means that the jid is a member of the room
	IsMember() bool
	// IsModerator returns true if this specific affiliation means that the jid is a moderator of the room
	IsModerator() bool
	// Name returns the string name of the affiliation type
	Name() string
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

func (*noneAffiliation) Name() string    { return AffiliationNone }
func (*outcastAffiliation) Name() string { return AffiliationOutcast }
func (*memberAffiliation) Name() string  { return AffiliationMember }
func (*adminAffiliation) Name() string   { return AffiliationAdmin }
func (*ownerAffiliation) Name() string   { return AffiliationOwner }

// AffiliationFromString returns an Affiliation from the given string, or an error if the string doesn't match a known affiliation type
func AffiliationFromString(s string) (Affiliation, error) {
	switch s {
	case AffiliationNone:
		return &noneAffiliation{}, nil
	case AffiliationOutcast:
		return &outcastAffiliation{}, nil
	case AffiliationMember:
		return &memberAffiliation{}, nil
	case AffiliationAdmin:
		return &adminAffiliation{}, nil
	case AffiliationOwner:
		return &ownerAffiliation{}, nil
	default:
		return nil, fmt.Errorf("unknown affiliation string: '%s'", s)
	}
}
