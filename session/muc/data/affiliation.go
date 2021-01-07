package data

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
	// IsAdmin will return true if this specific affiliation can modify persistent information
	IsAdmin() bool
	// IsBanned will return true if this specific affiliation means that the jid is banned from the room
	IsBanned() bool
	// IsMember returns true if this specific affiliation means that the jid is a member of the room
	IsMember() bool
	// IsModerator returns true if this specific affiliation means that the jid is a moderator of the room
	IsModerator() bool
	// IsOwner returns ture if this specific affiliation means that the jid is an owner of the room
	IsOwner() bool
	// Name returns the string name of the affiliation type
	Name() string
}

// NoneAffiliation is a representation of MUC's "none" affiliation
type NoneAffiliation struct{}

// OutcastAffiliation is a representation of MUC's "banned" affiliation
type OutcastAffiliation struct{}

// MemberAffiliation is a representation of MUC's "member" affiliation
type MemberAffiliation struct{}

// AdminAffiliation is a representation of MUC's "admin" affiliation
type AdminAffiliation struct{}

// OwnerAffiliation is a representation of MUC's "owner" affiliation
type OwnerAffiliation struct{}

// IsAdmin implements Affiliation interface
func (*NoneAffiliation) IsAdmin() bool { return false }

// IsAdmin implements Affiliation interface
func (*OutcastAffiliation) IsAdmin() bool { return false }

// IsAdmin implements Affiliation interface
func (*MemberAffiliation) IsAdmin() bool { return false }

// IsAdmin implements Affiliation interface
func (*AdminAffiliation) IsAdmin() bool { return true }

// IsAdmin implements Affiliation interface
func (*OwnerAffiliation) IsAdmin() bool { return false }

// IsBanned implements Affiliation interface
func (*NoneAffiliation) IsBanned() bool { return false }

// IsBanned implements Affiliation interface
func (*OutcastAffiliation) IsBanned() bool { return true }

// IsBanned implements Affiliation interface
func (*MemberAffiliation) IsBanned() bool { return false }

// IsBanned implements Affiliation interface
func (*AdminAffiliation) IsBanned() bool { return false }

// IsBanned implements Affiliation interface
func (*OwnerAffiliation) IsBanned() bool { return false }

// IsMember implements Affiliation interface
func (*NoneAffiliation) IsMember() bool { return false }

// IsMember implements Affiliation interface
func (*OutcastAffiliation) IsMember() bool { return false }

// IsMember implements Affiliation interface
func (*MemberAffiliation) IsMember() bool { return true }

// IsMember implements Affiliation interface
func (*AdminAffiliation) IsMember() bool { return true }

// IsMember implements Affiliation interface
func (*OwnerAffiliation) IsMember() bool { return true }

// IsModerator implements Affiliation interface
func (*NoneAffiliation) IsModerator() bool { return false }

// IsModerator implements Affiliation interface
func (*OutcastAffiliation) IsModerator() bool { return false }

// IsModerator implements Affiliation interface
func (*MemberAffiliation) IsModerator() bool { return false }

// IsModerator implements Affiliation interface
func (*AdminAffiliation) IsModerator() bool { return true }

// IsModerator implements Affiliation interface
func (*OwnerAffiliation) IsModerator() bool { return true }

// IsOwner implements Affiliation interface
func (*NoneAffiliation) IsOwner() bool { return false }

// IsOwner implements Affiliation interface
func (*OutcastAffiliation) IsOwner() bool { return false }

// IsOwner implements Affiliation interface
func (*MemberAffiliation) IsOwner() bool { return false }

// IsOwner implements Affiliation interface
func (*AdminAffiliation) IsOwner() bool { return false }

// IsOwner implements Affiliation interface
func (*OwnerAffiliation) IsOwner() bool { return true }

// Name implements Affiliation interface
func (*NoneAffiliation) Name() string { return AffiliationNone }

// Name implements Affiliation interface
func (*OutcastAffiliation) Name() string { return AffiliationOutcast }

// Name implements Affiliation interface
func (*MemberAffiliation) Name() string { return AffiliationMember }

// Name implements Affiliation interface
func (*AdminAffiliation) Name() string { return AffiliationAdmin }

// Name implements Affiliation interface
func (*OwnerAffiliation) Name() string { return AffiliationOwner }

// AffiliationFromString returns an Affiliation from the given string, or an error if the string doesn't match a known affiliation type
func AffiliationFromString(s string) (Affiliation, error) {
	switch s {
	case AffiliationNone:
		return &NoneAffiliation{}, nil
	case AffiliationOutcast:
		return &OutcastAffiliation{}, nil
	case AffiliationMember:
		return &MemberAffiliation{}, nil
	case AffiliationAdmin:
		return &AdminAffiliation{}, nil
	case AffiliationOwner:
		return &OwnerAffiliation{}, nil
	default:
		return nil, fmt.Errorf("unknown affiliation string: '%s'", s)
	}
}
