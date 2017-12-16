package data

import "encoding/xml"

// ClientIQ contains a specific information query request
type ClientIQ struct { // info/query
	XMLName xml.Name    `xml:"jabber:client iq"`
	From    string      `xml:"from,attr"`
	ID      string      `xml:"id,attr"`
	To      string      `xml:"to,attr"`
	Type    string      `xml:"type,attr"` // error, get, result, set
	Error   StanzaError `xml:"jabber:client error"`
	Bind    BindBind    `xml:"bind"`
	Query   []byte      `xml:",innerxml"`
}

// An EmptyReply results is in no XML.
type EmptyReply struct {
}
