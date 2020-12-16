package data

import "github.com/coyim/coyim/xmpp/data"

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

// OptionList description
type OptionList struct {
	label string
	value string
}

// NewOptionList description
func NewOptionList(l, v string) OptionList {
	return OptionList{l, v}
}

// ListSingleField description
type ListSingleField struct {
	value   string
	options []OptionList
}

// NewListSingleField description
func NewListSingleField(value string, options []data.FormFieldOptionX) *ListSingleField {
	return &ListSingleField{
		value:   value,
		options: populateSingleListOptions(options),
	}
}

// ListMultiField description
type ListMultiField struct {
	values  []string
	options []OptionList
}

// NewListMultiField description
func NewListMultiField(values []string, options []data.FormFieldOptionX) *ListMultiField {
	return &ListMultiField{
		values:  values,
		options: populateSingleListOptions(options),
	}
}

func populateSingleListOptions(options []data.FormFieldOptionX) (listOptions []OptionList) {
	for _, o := range options {
		listOptions = append(listOptions, NewOptionList(o.Label, o.Value))
	}
	return listOptions
}
