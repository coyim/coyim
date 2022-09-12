package data

import "encoding/xml"

// InBandRegistration represents an inband registration
type InBandRegistration struct {
	XMLName xml.Name `xml:"http://jabber.org/features/iq-register register,omitempty"`
}

// RegisterQuery contains register query information for creating a new account
type RegisterQuery struct {
	XMLName      xml.Name  `xml:"jabber:iq:register query"`
	Username     *xml.Name `xml:"username"`
	Password     *xml.Name `xml:"password"`
	Form         Form      `xml:"jabber:x:data x"`
	Instructions string    `xml:"instructions,omitempty"`
	Link         *OobLink  `xml:"jabber:x:oob x"`
	Datas        []BobData `xml:"data"`
}
