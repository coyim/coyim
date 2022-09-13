package data

import "encoding/xml"

// OobLink is a data element from http://xmpp.org/extensions/xep-0066.html.
type OobLink struct {
	XMLName     xml.Name `xml:"jabber:x:oob x"`
	URL         string   `xml:"url"`
	Description string   `xml:"desc"`
}
