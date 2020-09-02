package muc

import "fmt"

const (
	// RoleNone represents XMPP muc 'none' role
	RoleNone = "none"
	// RoleVisitor represents XMPP muc 'visitor' role
	RoleVisitor = "visitor"
	// RoleParticipant represents XMPP muc 'participant' role
	RoleParticipant = "participant"
	// RoleModerator represents XMPP muc 'moderator' role
	RoleModerator = "moderator"
)

// Role represents the specific role that a user has inside a specific room
type Role interface {
	// HasVoice returns true if the user can speak in this room
	HasVoice() bool
	// WithVoice returns the closest role upwards that has voice privilege. For Participants and Moderators, it returns itself, otherwise it returns participant
	WithVoice() Role
	// AsModerator returns the closest role upwards that can act as a moderator
	AsModerator() Role
	// Name returns the string name of the role type
	Name() string
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

func (*noneRole) Name() string        { return RoleNone }
func (*visitorRole) Name() string     { return RoleVisitor }
func (*participantRole) Name() string { return RoleParticipant }
func (*moderatorRole) Name() string   { return RoleModerator }

// RoleFromString returns the role object that matches the string given, or an error if the string given doesn't match a known role
func RoleFromString(s string) (Role, error) {
	switch s {
	case RoleNone:
		return &noneRole{}, nil
	case RoleVisitor:
		return &visitorRole{}, nil
	case RoleParticipant:
		return &participantRole{}, nil
	case RoleModerator:
		return &moderatorRole{}, nil
	default:
		return nil, fmt.Errorf("unknown role string: '%s'", s)
	}
}
