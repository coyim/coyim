package data

import (
	"encoding/xml"
)

// MUCOwnerQuery contains the deserialized information about a room configuration query
type MUCOwnerQuery struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/muc#owner query"`
	Form    *Form    `xml:",omitempty"`
}

// MUC contains information related with Presence x tag
type MUC struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/muc x"`
}

// MUCUserItem contains information related to role and affiliation
type MUCUserItem struct {
	XMLName     xml.Name `xml:"http://jabber.org/protocol/muc#user item"`
	Role        string   `xml:"role,attr,omitempty"` //moderator, participant, visitor
	Jid         string   `xml:"jid,attr,omitempty"`
	Affiliation string   `xml:"affiliation,attr,omitempty"` //owner, admin, member, none
}

// MUCUserStatus contains information related to status of the occupant
type MUCUserStatus struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/muc#user status"`
	Code    string   `xml:"code,attr,omitempty"`
}

// MUCUser contains information related to extended presence information about roles and affiliation
type MUCUser struct {
	XMLName xml.Name        `xml:"http://jabber.org/protocol/muc#user x"`
	Item    MUCUserItem     `xml:"item,omitempty"`
	Status  []MUCUserStatus `xml:"status,omitempty"`
}
