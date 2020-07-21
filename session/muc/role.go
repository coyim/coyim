package muc

import "fmt"

// Role represents the specific role that a user has inside a specific room
type Role interface {
	// HasVoice returns true if the user can speak in this room
	HasVoice() bool
	// WithVoice returns the closest role upwards that has voice privilege. For Participants and Moderators, it returns itself, otherwise it returns participant
	WithVoice() Role
	// AsModerator returns the closest role upwards that can act as a moderator
	AsModerator() Role
}

type noneRole struct{}
type visitorRole struct{}
type participantRole struct{}
type moderatorRole struct{}

func (*noneRole) HasVoice() bool        { return false }
func (*visitorRole) HasVoice() bool     { return false }
func (*participantRole) HasVoice() bool { return true }
func (*moderatorRole) HasVoice() bool   { return true }

func (*noneRole) WithVoice() Role        { return &participantRole{} }
func (*visitorRole) WithVoice() Role     { return &participantRole{} }
func (*participantRole) WithVoice() Role { return &participantRole{} }
func (*moderatorRole) WithVoice() Role   { return &moderatorRole{} }

func (*noneRole) AsModerator() Role        { return &moderatorRole{} }
func (*visitorRole) AsModerator() Role     { return &moderatorRole{} }
func (*participantRole) AsModerator() Role { return &moderatorRole{} }
func (*moderatorRole) AsModerator() Role   { return &moderatorRole{} }

// RoleFromString returns the role object that matches the string given, or an error if the string given doesn't match a known role
func RoleFromString(s string) (Role, error) {
	switch s {
	case "none":
		return &noneRole{}, nil
	case "visitor":
		return &visitorRole{}, nil
	case "participant":
		return &participantRole{}, nil
	case "moderator":
		return &moderatorRole{}, nil
	default:
		return nil, fmt.Errorf("unknown role string: '%s'", s)
	}
}
