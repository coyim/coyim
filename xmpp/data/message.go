package data

import "encoding/xml"

// Message structure used to send live messages
type Message struct {
	XMLName xml.Name `xml:"jabber:client message"`
	From    string   `xml:"from,attr"`
	ID      string   `xml:"id,attr"`
	To      string   `xml:"to,attr"`
	Type    string   `xml:"type,attr"`
	Body    string   `xml:"body"`
}
