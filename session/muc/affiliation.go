package muc

import "fmt"

type Affiliation interface {
	IsBanned() bool
	IsMember() bool
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
