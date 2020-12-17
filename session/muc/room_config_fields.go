package muc

const (
	// RoomConfigOptionModerators represents the field option for "moderators"
	RoomConfigOptionModerators = "moderators"
	// RoomConfigOptionParticipants represents the field option for "participants"
	RoomConfigOptionParticipants = "participants"
	// RoomConfigOptionAnyone represents the field opion for "anyone"
	RoomConfigOptionAnyone = "anyone"
	// RoomConfigOptionModerator represents the field option for "moderator"
	RoomConfigOptionModerator = "moderator"
	// RoomConfigOptionParticipant represents the field option for "participant"
	RoomConfigOptionParticipant = "participant"
	// RoomConfigOptionVisitor represents the field option for "visitor"
	RoomConfigOptionVisitor = "visitor"
	// RoomConfigOptionNone represents the field option for "none"
	RoomConfigOptionNone = "none"
	// RoomConfigOption10 represents the field option for "10"
	RoomConfigOption10 = "10"
	// RoomConfigOption20 represents the field option for "20"
	RoomConfigOption20 = "20"
	// RoomConfigOption30 represents the field option for "30"
	RoomConfigOption30 = "30"
	// RoomConfigOption50 represents the field option for "50"
	RoomConfigOption50 = "50"
	// RoomConfigOption100 represents the field option for "100"
	RoomConfigOption100 = "100"
)

// ConfigListSingleField description
type ConfigListSingleField interface {
	// UpdateField will update the list fields with the given "value" and "options"
	UpdateField(string, []string)
	// UpdateValue updates the current field value
	UpdateValue(string)
	// Value returns the current list value
	CurrentValue() string
	// Options returns the field options
	Options() []string
}

type configListSingleField struct {
	value   string
	options []string
}

func newConfigListSingleField(o []string) ConfigListSingleField {
	return &configListSingleField{
		options: o,
	}
}

func (cf *configListSingleField) UpdateField(v string, o []string) {
	cf.UpdateValue(v)
	if len(o) != 0 {
		cf.options = o
	}
}

func (cf *configListSingleField) UpdateValue(v string) {
	cf.value = v
}

func (cf *configListSingleField) CurrentValue() string {
	return cf.value
}

func (cf *configListSingleField) Options() []string {
	return cf.options
}

// ConfigListMultiField description
type ConfigListMultiField interface {
	// UpdateField will update the list fields with the given "values" and "options"
	UpdateField([]string, []string)
}

type configListMultiField struct {
	values  []string
	options []string
}

func newConfigListMultiField(o []string) ConfigListMultiField {
	return &configListMultiField{
		options: o,
	}
}

func (cf *configListMultiField) UpdateField(v, o []string) {
	cf.values = v
	cf.options = o
}
