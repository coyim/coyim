package data

// RoomConfigType is used to represent a room configuration field
type RoomConfigType int

const (
	// RoomConfigSupportsVoiceRequests represents rooms "SupportsVoiceRequests" config
	RoomConfigSupportsVoiceRequests RoomConfigType = iota
	// RoomConfigAllowsRegistration represents rooms "AllowsRegistration" config
	RoomConfigAllowsRegistration
	// RoomConfigPersistent represents rooms "Persistent" config
	RoomConfigPersistent
	// RoomConfigModerated represents rooms "Moderated" config
	RoomConfigModerated
	// RoomConfigOpen represents rooms "Open" config
	RoomConfigOpen
	// RoomConfigPasswordProtected represents rooms "PasswordProtected" config
	RoomConfigPasswordProtected
	// RoomConfigPublic represents rooms "Public" config
	RoomConfigPublic
	// RoomConfigLanguage represents rooms "Language" config
	RoomConfigLanguage
	// RoomConfigOccupantsCanChangeSubject represents rooms "OccupantsCanChangeSubject" config
	RoomConfigOccupantsCanChangeSubject
	// RoomConfigTitle represents rooms "Title" config
	RoomConfigTitle
	// RoomConfigDescription represents rooms "Description" config
	RoomConfigDescription
	// RoomConfigMembersCanInvite represents rooms "OccupantsCanInvite" config
	RoomConfigMembersCanInvite
	// RoomConfigAllowPrivateMessages represents rooms "AllowPrivateMessages" config
	RoomConfigAllowPrivateMessages
	// RoomConfigLogged represents rooms "Logged" config
	RoomConfigLogged
	// RoomConfigMaxHistoryFetch represents the maximum number of history messages returned by Room
	RoomConfigMaxHistoryFetch
)
