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
