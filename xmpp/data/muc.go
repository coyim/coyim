package data

import (
	"encoding/xml"
)

// RoomConfigurationQuery contains the deserialized information about a room configuration query
type RoomConfigurationQuery struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/muc#owner query"`
	Form    *Form    `xml:",omitempty"`
}

//MUCExtra contains information related with Presence x tag
type MUCExtra struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/muc x"`
}

//MUCExtendedPresenceInfoItem contains information related to role and affiliation
type MUCExtendedPresenceInfoItem struct {
	XMLName     xml.Name `xml:"http://jabber.org/protocol/muc#user item"`
	Role        string   `xml:"role,attr,omitempty"` //moderator, participant, visitor
	Jid         string   `xml:"jid,attr,omitempty"`
	Affiliation string   `xml:"affiliation,attr,omitempty"` //owner, admin, member, none
}

//MUCExtendedPresenceInfoStatus contains information related to status of the occupant
type MUCExtendedPresenceInfoStatus struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/muc#user status"`
	Code    string   `xml:"code,attr,omitempty"`
}

//MUCExtendedPresenceInfo contains information related to extended presence information about roles and affiliation
type MUCExtendedPresenceInfo struct {
	XMLName xml.Name                        `xml:"http://jabber.org/protocol/muc#user x"`
	Item    MUCExtendedPresenceInfoItem     `xml:"item,omitempty"`
	Status  []MUCExtendedPresenceInfoStatus `xml:"status,omitempty"`
}
