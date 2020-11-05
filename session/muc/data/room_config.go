package data

// RoomConfigType is used to represent a room configuration field
type RoomConfigType int

const (
	// RoomConfigSupportsVoiceRequests represents room's "SupportsVoiceRequests" config
	RoomConfigSupportsVoiceRequests RoomConfigType = iota
	// RoomConfigAllowsRegistration represents room's "AllowsRegistration" config
	RoomConfigAllowsRegistration
	// RoomConfigPersistent represents room's "Persistent" config
	RoomConfigPersistent
	// RoomConfigModerated represents room's "Moderated" config
	RoomConfigModerated
	// RoomConfigOpen represents room's "Open" config
	RoomConfigOpen
	// RoomConfigPasswordProtected represents room's "PasswordProtected" config
	RoomConfigPasswordProtected
	// RoomConfigPublic represents room's "Public" config
	RoomConfigPublic
	// RoomConfigLanguage represents room's "Language" config
	RoomConfigLanguage
	// RoomConfigOccupantsCanChangeSubject represents room's "OccupantsCanChangeSubject" config
	RoomConfigOccupantsCanChangeSubject
	// RoomConfigTitle represents room's "Title" config
	RoomConfigTitle
	// RoomConfigDescription represents room's "Description" config
	RoomConfigDescription
	// RoomConfigMembersCanInvite represents room's "OccupantsCanInvite" config
	RoomConfigMembersCanInvite
	// RoomConfigAllowPrivateMessages represents room's "AllowPrivateMessages" config
	RoomConfigAllowPrivateMessages
	// RoomConfigLogged represents room's "Logged" config
	RoomConfigLogged
	// RoomConfigMaxHistoryFetch represents the maximum number of history messages returned by Room
	RoomConfigMaxHistoryFetch
)

// RoomConfig represents the room's configuration values
type RoomConfig struct {
	SupportsVoiceRequests     bool
	AllowsRegistration        bool
	Anonymity                 string
	Persistent                bool
	Moderated                 bool
	Open                      bool
	PasswordProtected         bool
	Public                    bool
	Language                  string
	OccupantsCanChangeSubject bool
	Title                     string
	Description               string
	Occupants                 int
	MembersCanInvite          bool
	OccupantsCanInvite        bool
	AllowPrivateMessages      string // This can be 'anyone', 'participants', 'moderators', 'none'
	ContactJid                string
	Logged                    bool // Notice that this will not always be correct for all servers
	MaxHistoryFetch           int
}
