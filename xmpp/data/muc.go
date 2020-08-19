package data

import (
	"encoding/xml"
)

// MUC contains information related with Presence x tag
type MUC struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/muc x"`
}

// MUCStatus contains information related to status of the presence or message
type MUCStatus struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/muc status"`
	Code    string   `xml:"code,attr,omitempty"`
}

// MUCUser contains information related to extended presence information about roles and affiliation
type MUCUser struct {
	XMLName xml.Name        `xml:"http://jabber.org/protocol/muc#user x"`
	Item    *MUCUserItem    `xml:"item,omitempty"`
	Status  []MUCUserStatus `xml:"status,omitempty"`
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

// MUCRoomConfiguration contains the deserialized information about a room configuration query
type MUCRoomConfiguration struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/muc#owner query"`
	Form    *Form    `xml:",omitempty"`
}

// MUCNotAuthorized inform user that a password is required
type MUCNotAuthorized struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-stanzas not-authorized,omitempty"`
}

// MUCForbidden inform user that he or she is banned from the room
type MUCForbidden struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-stanzas forbidden,omitempty"`
}

// MUCItemNotFound inform user that the room does not exist
type MUCItemNotFound struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-stanzas item-not-found,omitempty"`
}

// MUCNotAllowed inform user that room creation is restricted
type MUCNotAllowed struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-stanzas not-allowed,omitempty"`
}

// MUCNotAcceptable inform user that the reserved roomnick must be used
type MUCNotAcceptable struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-stanzas not-acceptable,omitempty"`
}

// MUCRegistrationRequired inform user that he or she is not on the member list
type MUCRegistrationRequired struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-stanzas registration-required,omitempty"`
}

// MUCConflict inform user that his or her desired room nickname is in use or
// registered by another user
type MUCConflict struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-stanzas conflict,omitempty"`
}

// MUCServiceUnavailable inform user that the maximum number of users has been reached
type MUCServiceUnavailable struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-stanzas service-unavailable,omitempty"`
}
