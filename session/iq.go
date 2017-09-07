package session

import (
	"bytes"
	"fmt"

	"github.com/twstrike/coyim/xmpp"
	"github.com/twstrike/coyim/xmpp/data"
)

func (s *session) sendIQError(stanza *data.ClientIQ, reply interface{}) {
	s.sendIQReply(stanza, "error", reply)
}

func (s *session) sendIQResult(stanza *data.ClientIQ, reply interface{}) {
	s.sendIQReply(stanza, "result", reply)
}

func (s *session) sendIQReply(stanza *data.ClientIQ, tp string, reply interface{}) {
	if err := s.conn.SendIQReply(stanza.From, tp, stanza.ID, reply); err != nil {
		s.alert("Failed to send IQ message: " + err.Error())
	}
}

func discoIQ(s *session, _ *data.ClientIQ) (ret interface{}, ignore bool) {
	// TODO: We should ensure that there is no "node" entity on this query, since we don't support that.
	// In the case of a "node", we should return  <service-unavailable/>
	s.info("IQ: http://jabber.org/protocol/disco#info query")
	return xmpp.DiscoveryReply(s.GetConfig().Account), false
}

func versionIQ(s *session, _ *data.ClientIQ) (ret interface{}, ignore bool) {
	s.info("IQ: jabber:iq:version query")
	return s.receivedIQVersion(), false
}

func rosterIQ(s *session, stanza *data.ClientIQ) (ret interface{}, ignore bool) {
	s.info("IQ: jabber:iq:roster query")
	return s.receivedIQRosterQuery(stanza)
}

func unknownIQ(s *session, stanza *data.ClientIQ) (ret interface{}, ignore bool) {
	s.info(fmt.Sprintf("Unknown IQ: %s", bytes.NewBuffer(stanza.Query)))
	return nil, false
}

type iqFunction func(*session, *data.ClientIQ) (interface{}, bool)

var knownIQs = map[string]iqFunction{}

func registerKnownIQ(stanzaType, fullName string, f iqFunction) {
	knownIQs[stanzaType+" "+fullName] = f
}

func getIQHandler(stanzaType, namespace, local string) iqFunction {
	f, ok := knownIQs[fmt.Sprintf("%s %s %s", stanzaType, namespace, local)]
	if ok {
		return f
	}
	return unknownIQ
}

func init() {
	registerKnownIQ("get", "http://jabber.org/protocol/disco#info query", discoIQ)
	registerKnownIQ("get", "jabber:iq:version query", versionIQ)
	registerKnownIQ("set", "jabber:iq:roster query", rosterIQ)
	registerKnownIQ("result", "jabber:iq:roster query", rosterIQ)
}

func (s *session) processIQ(stanza *data.ClientIQ) (ret interface{}, ignore bool) {
	if nspace, local, ok := tryDecodeXML(stanza.Query); ok {
		return getIQHandler(stanza.Type, nspace, local)(s, stanza)
	}
	return nil, false
}
