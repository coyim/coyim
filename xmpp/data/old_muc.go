package data

import (
	"encoding/xml"
	"fmt"
)

//See: Section 4.1

// LegacyOldDoNotUseRoom represents a chat room
type LegacyOldDoNotUseRoom struct {
	ID, Service string
}

// JID returns the JID for this room
func (o *LegacyOldDoNotUseRoom) JID() string {
	return fmt.Sprintf("%s@%s", o.ID, o.Service)
}

//See: Section 4.1

// LegacyOldDoNotUseOccupant represents a person in a chat room
type LegacyOldDoNotUseOccupant struct {
	LegacyOldDoNotUseRoom
	Handle string
}

// JID returns the JID for this occupant
func (o *LegacyOldDoNotUseOccupant) JID() string {
	return fmt.Sprintf("%s/%s", o.LegacyOldDoNotUseRoom.JID(), o.Handle)
}

// LegacyOldDoNotUseRoomConfigurationQuery contains the deserialized information about a room configuration query
// See: Section "10.2 Subsequent Room Configuration"
type LegacyOldDoNotUseRoomConfigurationQuery struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/muc#owner query"`
	Form    *Form    `xml:",omitempty"`
}

// See: Section 15.5.4 muc#roominfo FORM_TYPE

// LegacyOldDoNotUseRoomInfoForm contains the information about a room
type LegacyOldDoNotUseRoomInfoForm struct {
	MaxHistoryFetch string   `form-field:"muc#maxhistoryfetch"`
	ContactJID      []string `form-field:"muc#roominfo_contactjid"`
	Description     string   `form-field:"muc#roominfo_description"`
	Language        string   `form-field:"muc#roominfo_language"`
	LDAPGroup       string   `form-field:"muc#roominfo_ldapgroup"`
	Logs            string   `form-field:"muc#roominfo_logs"`
	Occupants       int      `form-field:"muc#roominfo_occupants"`
	Subject         string   `form-field:"muc#roominfo_subject"`
	SubjectMod      bool     `form-field:"muc#roominfo_subjectmod"`
}

//See: Section 4.2

// LegacyOldDoNotUseRoomType contains information about the different options for the room
type LegacyOldDoNotUseRoomType struct {
	Public bool
	//vs Hidden bool

	Open bool
	//vs MembersOnly bool

	Moderated bool
	//vs Unmoderated bool

	SemiAnonymous bool
	//vs NonAnonymous bool

	PasswordProtected bool
	//vs Unsecured bool

	Persistent bool
	//vs Temporary bool
}

//TODO: Ahh, naming

// LegacyOldDoNotUseRoomInfo contains room information
type LegacyOldDoNotUseRoomInfo struct {
	LegacyOldDoNotUseRoomInfoForm `form-type:"http://jabber.org/protocol/muc#roominfo"`
	LegacyOldDoNotUseRoomType
}
