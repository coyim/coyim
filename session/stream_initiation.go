package session

import (
	"bytes"
	"encoding/xml"
	"fmt"

	"github.com/twstrike/coyim/xmpp/data"
)

func init() {
	registerKnownIQ("set", "http://jabber.org/protocol/si si", streamInitIQ)
}

var supportedSIProfiles = map[string]func(*session, *data.ClientIQ, data.SI) (interface{}, bool){
	"http://jabber.org/protocol/si/profile/file-transfer": fileStreamInitIQ,
}

func streamInitIQ(s *session, stanza *data.ClientIQ) (ret interface{}, ignore bool) {
	s.info("IQ: stream initiation")
	var si data.SI
	if err := xml.NewDecoder(bytes.NewBuffer(stanza.Query)).Decode(&si); err != nil {
		s.warn(fmt.Sprintf("Failed to parse stream initiation: %v", err))
		return nil, false
	}

	prof, ok := supportedSIProfiles[si.Profile]
	if !ok {
		s.warn(fmt.Sprintf("Unsupported SI profile: %v", si.Profile))
		return nil, false
	}

	return prof(s, stanza, si)
}
