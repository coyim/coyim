package data

import "encoding/xml"

// EncryptedData is an element from an unpublished potential XEP
type EncryptedData struct {
	XMLName              xml.Name              `xml:"http://jabber.org/protocol/si/profile/encrypted-data-transfer data"`
	Type                 string                `xml:"type,omitempty"`
	MediaType            string                `xml:"media-type,omitempty"`
	Name                 string                `xml:"name,omitempty"`
	Size                 int64                 `xml:"size,omitempty"`
	Desc                 string                `xml:"desc,omitempty"`
	Date                 string                `xml:"date,omitempty"`
	Range                *FileRange            `xml:",omitempty"`
	Hash                 *Hash                 `xml:"hash,omitempty"`
	EncryptionParameters *EncryptionParameters `xml:"encryption,omitempty"`
}

// EncryptionParameters is an element from an unpublished potential XEP
type EncryptionParameters struct {
	XMLName       xml.Name                `xml:"http://jabber.org/protocol/si/profile/encrypted-data-transfer encryption"`
	Type          string                  `xml:"type,attr"`
	IV            string                  `xml:"iv,attr"`
	MAC           string                  `xml:"mac,attr"`
	EncryptionKey *EncryptionKeyParameter `xml:"encryption-key"`
	MACKey        *MACKeyParameter        `xml:"mac-key"`
}

// EncryptionKeyParameter is an element from an unpublished potential XEP
type EncryptionKeyParameter struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/si/profile/encrypted-data-transfer encryption-key"`
	Type    string   `xml:"type,attr"` // 'static' or 'external'
	Value   string   `xml:"value,attr,omitempty"`
}

// MACKeyParameter is an element from an unpublished potential XEP
type MACKeyParameter struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/si/profile/encrypted-data-transfer mac-key"`
	Type    string   `xml:"type,attr"` // 'static' or 'external'
	Value   string   `xml:"value,attr,omitempty"`
}
