package data

import (
	"bytes"
	"encoding/xml"
	"errors"
)

// VCard contains vcard
type VCard struct {
	XMLName xml.Name `xml:"vcard-temp vCard"`

	FullName string `xml:"FN"`
	Nickname string `xml:"NICKNAME"`
}

// ParseVCard extracts vcard information from the given Stanza.
func ParseVCard(reply Stanza) (VCard, error) {
	iq, ok := reply.Value.(*ClientIQ)
	if !ok {
		return VCard{}, errors.New("xmpp: vcard request resulted in tag of type " + reply.Name.Local)
	}

	var vcard VCard
	if err := xml.NewDecoder(bytes.NewBuffer(iq.Query)).Decode(&vcard); err != nil {
		return VCard{}, err
	}
	return vcard, nil
}
