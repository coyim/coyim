package data

import (
	"encoding/xml"
	"fmt"
)

//See: Section 4.1

// Room represents a chat room
type Room struct {
	ID, Service string
}

// JID returns the JID for this room
func (o *Room) JID() string {
	return fmt.Sprintf("%s@%s", o.ID, o.Service)
}

//See: Section 4.1

// Occupant represents a person in a chat room
type Occupant struct {
	Room
	Handle string
}

// JID returns the JID for this occupant
func (o *Occupant) JID() string {
	return fmt.Sprintf("%s/%s", o.Room.JID(), o.Handle)
}

// RoomConfigurationQuery contains the deserialized information about a room configuration query
// See: Section "10.2 Subsequent Room Configuration"
type RoomConfigurationQuery struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/muc#owner query"`
	Form    *Form    `xml:,omitempty`
}

// See: Section 15.5.4 muc#roominfo FORM_TYPE

// RoomInfoForm contains the information about a room
type RoomInfoForm struct {
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

// RoomType contains information about the different options for the room
type RoomType struct {
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

// RoomInfo contains room information
type RoomInfo struct {
	RoomInfoForm `form-type:"http://jabber.org/protocol/muc#roominfo"`
	RoomType
}
