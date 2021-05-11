package data

// RoomConfigType is used to represent a room configuration field
type RoomConfigType int

const (
	// RoomConfigSupportsVoiceRequests represents the field "SupportsVoiceRequests" of the room configuration info
	RoomConfigSupportsVoiceRequests RoomConfigType = iota
	// RoomConfigAllowsRegistration represents the field "AllowsRegistration" of the room configuration info
	RoomConfigAllowsRegistration
	// RoomConfigPersistent represents the field "Persistent" of the room configuration info
	RoomConfigPersistent
	// RoomConfigModerated represents the field "Moderated" of the room configuration info
	RoomConfigModerated
	// RoomConfigOpen represents the field "Open" of the room configuration info
	RoomConfigOpen
	// RoomConfigPasswordProtected represents the field "PasswordProtected" of the room configuration info
	RoomConfigPasswordProtected
	// RoomConfigPublic represents the field "Public" of the room configuration info
	RoomConfigPublic
	// RoomConfigLanguage represents the field "Language" of the room configuration info
	RoomConfigLanguage
	// RoomConfigOccupantsCanChangeSubject represents the field "OccupantsCanChangeSubject" of the room configuration info
	RoomConfigOccupantsCanChangeSubject
	// RoomConfigTitle represents the field "Title" of the room configuration info
	RoomConfigTitle
	// RoomConfigDescription represents the field "Description" of the room configuration info
	RoomConfigDescription
	// RoomConfigMembersCanInvite represents the field "OccupantsCanInvite" of the room configuration info
	RoomConfigMembersCanInvite
	// RoomConfigAllowPrivateMessages represents the field "AllowPrivateMessages" of the room configuration info
	RoomConfigAllowPrivateMessages
	// RoomConfigLogged represents the field "Logged" of the room configuration info
	RoomConfigLogged
	// RoomConfigMaxHistoryFetch represents the field "MaxHistoryFetch" of the room configuration info
	RoomConfigMaxHistoryFetch
)
