package data

// RoomConfigType is used to represent a room configuration field
type RoomConfigType int

const (
	// RoomConfigSupportsVoiceRequests represents the "SupportsVoiceRequests" configuration info of the room
	RoomConfigSupportsVoiceRequests RoomConfigType = iota
	// RoomConfigAllowsRegistration represents the "AllowsRegistration" configuration info of the room
	RoomConfigAllowsRegistration
	// RoomConfigPersistent represents the "Persistent" configuration info of the room
	RoomConfigPersistent
	// RoomConfigModerated represents the "Moderated" configuration info of the room
	RoomConfigModerated
	// RoomConfigOpen represents the "Open" configuration info of the room
	RoomConfigOpen
	// RoomConfigPasswordProtected represents the "PasswordProtected" configuration info of the room
	RoomConfigPasswordProtected
	// RoomConfigPublic represents the "Public" configuration info of the room
	RoomConfigPublic
	// RoomConfigLanguage represents the "Language" configuration info of the room
	RoomConfigLanguage
	// RoomConfigOccupantsCanChangeSubject represents the "OccupantsCanChangeSubject" configuration info of the room
	RoomConfigOccupantsCanChangeSubject
	// RoomConfigTitle represents the "Title" configuration info of the room
	RoomConfigTitle
	// RoomConfigDescription represents the "Description" configuration info of the room
	RoomConfigDescription
	// RoomConfigMembersCanInvite represents the "OccupantsCanInvite" configuration info of the room
	RoomConfigMembersCanInvite
	// RoomConfigAllowPrivateMessages represents the "AllowPrivateMessages" configuration info of the room
	RoomConfigAllowPrivateMessages
	// RoomConfigLogged represents the "Logged" configuration info of the room
	RoomConfigLogged
	// RoomConfigMaxHistoryFetch represents the "MaxHistoryFetch" configuration info of the room
	RoomConfigMaxHistoryFetch
)
