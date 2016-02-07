package data

import "encoding/xml"

// VersionQuery represents a version query
type VersionQuery struct {
	XMLName xml.Name `xml:"jabber:iq:version query"`
}

// VersionReply contains a version reply
type VersionReply struct {
	XMLName xml.Name `xml:"jabber:iq:version query"`
	Name    string   `xml:"name"`
	Version string   `xml:"version"`
	OS      string   `xml:"os"`
}
