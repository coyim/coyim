package data

// RoomInfo represents the room's information form values
// (https://xmpp.org/extensions/xep-0045.html#registrar-formtype-roominfo)
type RoomInfo struct {
	// muc#maxhistoryfetch
	// Maximum Number of History Messages Returned by Room
	MaxHistoryFetch int

	// muc#roominfo_contactjid
	// Contact Addresses (normally, room owner or owners)
	ContactJid string

	// muc#roominfo_description
	// Short Description of Room
	Description string

	//muc#roominfo_lang
	// Natural Language for Room Discussions
	Language string

	// muc#roominfo_ldapgroup
	// An associated LDAP group that defines room membership; this should be an LDAP
	// Distinguished Name according to an implementation-specific or
	// deployment-specific definition of a group.
	LdapGroup string

	// muc#roominfo_logs
	// URL for Archived Discussion Logs
	Logs string

	// muc#roominfo_occupants
	// Current Number of Occupants in Room
	Occupants int

	// muc#roominfo_subject
	// Current Discussion Topic
	Topic string

	// muc#roominfo_subjectmod
	// The room subject can be modified by participants
	OccupantsCanChangeSubject bool
}
