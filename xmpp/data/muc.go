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

//ExtendedPresenceInfoItem contains information related to role and affiliation
type ExtendedPresenceInfoItem struct {
	XMLName     xml.Name
	Role        string `xml:"role,attr,omitempty"` //moderator, participant, visitor
	Jid         string `xml:"jid,attr,omitempty"`
	Affiliation string `xml:"affiliation,attr,omitempty"` //owner, admin, member, none
}

//ExtendedPresenceInfoStatus contains information related to status of the occupant
type ExtendedPresenceInfoStatus struct {
	XMLName xml.Name
	Code    string `xml:"code,attr,omitempty"`
}

//ExtendedPresenceInfo contains information related to extended presence information about roles and affiliation
type ExtendedPresenceInfo struct {
	XMLName xml.Name                     `xml:"http://jabber.org/protocol/muc#user x"`
	Item    ExtendedPresenceInfoItem     `xml:"item,omitempty"`
	Status  []ExtendedPresenceInfoStatus `xml:"status,omitempty"`
}
