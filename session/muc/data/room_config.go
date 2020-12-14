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

// RoomConfig represents the room configuration form values
// (https://xmpp.org/extensions/xep-0045.html#registrar-formtype-owner)
type RoomConfig struct {
	// muc#maxhistoryfetch
	// Maximum Number of History Messages Returned by Room
	MaxHistoryFetch roomConfigTextSingleField

	// muc#roomconfig_allowpm
	// Roles that May Send Private Messages
	AllowPrivateMessages roomConfigListSingleField

	// muc#roomconfig_allowinvites
	// Whether to Allow Occupants to Invite Others
	OccupantsCanInvite roomConfigBooleanField

	// muc#roomconfig_changesubject
	// Whether to Allow Occupants to Change Subject
	OccupantsCanChangeSubject roomConfigBooleanField

	// muc#roomconfig_enablelogging
	// Whether to Enable Public Logging of Room Conversations
	Logged roomConfigBooleanField

	// muc#roomconfig_getmemberlist
	// Roles and Affiliations that May Retrieve Member List
	RetrieveMembersList roomConfigListMultiField

	// muc#roomconfig_lang
	// Natural Language for Room Discussions
	Language roomConfigTextSingleField

	// muc#roomconfig_pubsub
	// XMPP URI of Associated Publish-Subscribe Node
	AssociatedPublishSubscribeNode roomConfigTextSingleField

	// muc#roomconfig_maxusers
	// Maximum Number of Room Occupants
	MaxOccupantsNumber roomConfigListSingleField

	// muc#roomconfig_membersonly
	// Whether to Make Room Members-Only
	MembersOnly roomConfigBooleanField

	// muc#roomconfig_moderatedroom
	// Whether to Make Room Moderated
	Moderated roomConfigBooleanField

	// muc#roomconfig_passwordprotectedroom
	// Whether a Password is Required to Enter
	PasswordProtected roomConfigBooleanField

	// muc#roomconfig_persistentroom
	// Whether to Make Room Persistent
	Persistent roomConfigBooleanField

	// muc#roomconfig_presencebroadcast
	// Roles for which Presence is Broadcasted
	PresenceBroadcast roomConfigListMultiField

	// muc#roomconfig_publicroom
	// Whether to Allow Public Searching for Room
	Public roomConfigBooleanField

	// muc#roomconfig_roomadmins
	// Full List of Room Admins
	Admins roomConfigJidMultiField

	// muc#roomconfig_roomdesc
	// Short Description of Room
	Description roomConfigTextSingleField

	// muc#roomconfig_roomname
	// Natural-Language Room Name
	Title roomConfigTextSingleField

	// muc#roomconfig_roomowners
	// Full List of Room Owners
	Owners roomConfigJidMultiField

	// muc#roomconfig_roomsecret
	// The Room Password
	Password roomConfigTextSingleField

	// muc#roomconfig_whois
	// Affiliations that May Discover Real JIDs of Occupants
	Whois roomConfigListSingleField
}
