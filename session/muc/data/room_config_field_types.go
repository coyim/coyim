package data

import "github.com/coyim/coyim/xmpp/jid"

const (
	// PrivateMessageAnyoneRole description
	PrivateMessageAnyoneRole = "anyone"
	// PrivateMessageParticipantsRole description
	PrivateMessageParticipantsRole = "participants"
	// PrivateMessageModeratorsRole description
	PrivateMessageModeratorsRole = "moderators"
	// PrivateMessageNoneRole description
	PrivateMessageNoneRole = "none"
)

const (
	// MaximumOccupants10 description
	MaximumOccupants10 = "10"
	// MaximumOccupants20 description
	MaximumOccupants20 = "20"
	// MaximumOccupants30 description
	MaximumOccupants30 = "30"
	// MaximumOccupants50 description
	MaximumOccupants50 = "50"
	// MaximumOccupants100 description
	MaximumOccupants100 = "100"
	// MaximumOccupantsNone description
	MaximumOccupantsNone = "none"
)

const (
	// RoomConfigModeratorRole description
	RoomConfigModeratorRole = "moderator"
	// RoomConfigParticipantRole description
	RoomConfigParticipantRole = "participant"
	// RoomConfigVisitorRole description
	RoomConfigVisitorRole = "visitor"
)

type roomConfigTextSingleField struct {
	value string
}

func (f *roomConfigTextSingleField) SetValue(v string) {
	f.value = v
}

type roomConfigBooleanField struct {
	value bool
}

func (f *roomConfigBooleanField) SetValue(v bool) {
	f.value = v
}

type roomConfigListSingleField struct {
	value   string
	options []string
}

func (f *roomConfigListSingleField) SetValue(v string) {
	f.value = v
}

func (f *roomConfigListSingleField) SetOptions(o []string) {
	f.options = o
}

type roomConfigListMultiField struct {
	values  []string
	options []string
}

func (f *roomConfigListMultiField) SetValues(v []string) {
	f.values = v
}

func (f *roomConfigListMultiField) SetOptions(o []string) {
	f.options = o
}

type roomConfigPrivateField struct {
	value string
}

func (f *roomConfigPrivateField) SetValue(v string) {
	f.value = v
}

type roomConfigJidMultiField struct {
	values []jid.Any
}

func (f *roomConfigJidMultiField) SetValues(v []jid.Any) {
	f.values = v
}
